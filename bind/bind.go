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
			return &value.MapValue{Value: m}
		default:
			return &value.ErrorValue{ErrorMessage: fmt.Sprintf(
				"Cannot bind value of type #%T to complang yet: %#V", v, v)}
		}
	}
}
