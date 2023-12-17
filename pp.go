package complang

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

func pretty(value any, maxHeight, maxWidth int) string {
	value = unpack(value)
	var buf bytes.Buffer
	err := encodeYaml(value, &buf)
	if err != nil {
		buf.Reset()
		fmt.Fprintf(&buf, "%#v", value)
	}
	s := buf.String()
	parts := strings.Split(s, "\n")
	for i := 0; i < len(parts); i++ {
		if maxWidth >= 0 && len(parts[i]) > maxWidth {
			parts[i] = parts[i][0:maxWidth] + "..."
		}
	}
	if maxHeight >= 0 && len(parts) > maxHeight {
		return strings.Join(parts[0:32], "\n") + "\n..."
	}
	return s
}

func encodeYaml(v any, w io.Writer) (err error) {
	defer func() {
		if e := recover(); e != nil && err == nil {
			if e, ok := e.(error); ok {
				err = e
			}
		}
	}()
	return yaml.NewEncoder(w).Encode(v)
}

func unpack(v any) any {
	switch v := v.(type) {
	case MapValue:
		result := map[string]any{}
		for k, vv := range v {
			result[k] = unpack(vv)
		}
		return result
	case SliceValue:
		result := []any{}
		for _, e := range v {
			result = append(result, unpack(e))
		}
		return result
	default:
		return v
	}
}
