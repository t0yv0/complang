package bind

import (
	"fmt"
	"reflect"

	"github.com/t0yv0/complang/value"
)

func Value(v any) value.Value {
	switch v := v.(type) {
	case nil:
		return &value.NullValue{}
	case value.Value:
		return v
	case reflect.Value:
		return Value(v.Interface())
	case string:
		return &value.StringValue{Value: v}
	case bool:
		return &value.BoolValue{Value: v}
	default:
		vv := reflect.ValueOf(v)
		switch {
		case vv.Kind() == reflect.Slice || vv.Kind() == reflect.Array:
			vs := []value.Value{}
			for i := 0; i < vv.Len(); i++ {
				vs = append(vs, Value(vv.Index(i)))
			}
			return &value.SliceValue{Value: vs}
		case vv.Kind() == reflect.Map && vv.Type().Key().Kind() == reflect.String:
			m := map[value.Symbol]value.Value{}
			for _, key := range vv.MapKeys() {
				m[value.NewSymbol(key.String())] = Value(vv.MapIndex(key))
			}
			return &value.MapValue{Value: m}
		case vv.Kind() == reflect.Struct:
			m := map[value.Symbol]value.Value{}
			for i := 0; i < vv.NumField(); i++ {
				if !vv.Type().Field(i).IsExported() {
					continue
				}
				nn := value.NewSymbol(vv.Type().Field(i).Name)
				m[nn] = Value(vv.Field(i))
			}
			for i := 0; i < vv.Type().NumMethod(); i++ {
				me := vv.Type().Method(i)
				if !me.IsExported() {
					continue
				}
				m[value.NewSymbol(me.Name)] = bindMethod(vv, me)
			}
			return &value.MapValue{Value: m}

		default:
			return &value.ErrorValue{ErrorMessage: fmt.Sprintf(
				"Cannot bind value of type #%T to complang yet: %#V", v, v)}
		}
	}
}

func FromValue(v value.Value) (any, error) {
	switch v := v.(type) {
	case *value.StringValue:
		return v.Value, nil
	default:
		err := fmt.Errorf("Cannot bind complang value back to Go yet: %s", v.Show())
		return Value(err), nil
	}
}

func bindMethod(vv reflect.Value, me reflect.Method) value.Value {
	mh := vv.MethodByName(me.Name)
	params := []value.Symbol{}
	for i := 1; i < me.Type.NumIn(); i++ {
		params = append(params, value.NewSymbol(fmt.Sprintf("a%d", i)))
	}
	return &value.CustomValue{ValueLike: &value.Closure{
		Params: params,
		Call: func(env value.Env) value.Value {
			args := []reflect.Value{}
			for _, p := range params {
				pv, ok := env.Lookup(p)
				if !ok {
					return Value(fmt.Errorf("Unbound %v", p))
				}
				x, err := FromValue(pv)
				if err != nil {
					return Value(err)
				}
				args = append(args, reflect.ValueOf(x))
			}
			ret := mh.Call(args)
			if len(ret) == 1 {
				return Value(ret[0])
			} else {
				return &value.NullValue{}
			}
		},
	}}
}
