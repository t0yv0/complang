package complang

import (
	"context"
	"fmt"
	"reflect"
)

func BindValue(v any) Value {
	switch v := v.(type) {
	case nil:
		return NullValue{}
	case Value:
		return v
	case reflect.Value:
		return BindValue(v.Interface())
	case string:
		return StringValue{v}
	case bool:
		return BoolValue{v}
	default:
		vv := reflect.ValueOf(v)
		switch {
		case vv.Kind() == reflect.Slice || vv.Kind() == reflect.Array:
			vs := []Value{}
			for i := 0; i < vv.Len(); i++ {
				vs = append(vs, BindValue(vv.Index(i)))
			}
			return SliceValue(vs)
		case vv.Kind() == reflect.Map && vv.Type().Key().Kind() == reflect.String:
			m := map[string]Value{}
			for _, key := range vv.MapKeys() {
				m[key.String()] = BindValue(vv.MapIndex(key))
			}
			return MapValue(m)
		case vv.Kind() == reflect.Struct:
			m := map[string]Value{}
			for i := 0; i < vv.NumField(); i++ {
				if !vv.Type().Field(i).IsExported() {
					continue
				}
				nn := vv.Type().Field(i).Name
				m[nn] = BindValue(vv.Field(i))
			}
			for i := 0; i < vv.Type().NumMethod(); i++ {
				me := vv.Type().Method(i)
				if !me.IsExported() {
					continue
				}
				m[me.Name] = bindMethod(vv, me)
			}
			return MapValue(m)

		default:
			return Error{fmt.Sprintf(
				"Cannot bind value of type %T to complang yet: %#V", v, v)}
		}
	}
}

func UnbindValue(ctx context.Context, v Value) (any, error) {
	switch v := v.(type) {
	case StringValue:
		return v.Text, nil
	default:
		return nil, fmt.Errorf("Cannot bind complang value back to Go yet: %s",
			Show(ctx, v))
	}
}

func bindMethod(vv reflect.Value, me reflect.Method) Value {
	mh := vv.MethodByName(me.Name)
	params := []string{}
	for i := 1; i < me.Type.NumIn(); i++ {
		params = append(params, fmt.Sprintf("a%d", i))
	}
	c := func(ctx context.Context, env Env) Value {
		args := []reflect.Value{}
		for _, p := range params {
			pv, ok := env.Lookup(p)
			if !ok {
				return Error{fmt.Sprintf("Unbound %v", p)}
			}
			x, err := UnbindValue(ctx, pv)
			if err != nil {
				return BindValue(err)
			}
			args = append(args, reflect.ValueOf(x))
		}
		ret := mh.Call(args)
		if len(ret) == 1 {
			return BindValue(ret[0])
		} else {
			return NullValue{}
		}
	}
	return Closure{
		Params: params,
		Call:   c,
	}
}
