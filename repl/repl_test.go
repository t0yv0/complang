package repl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t0yv0/complang/value"
)

func TestReadEvalComplete(t *testing.T) {
	inter := &complangInterpreter{env: value.MapEnv{map[value.Symbol]value.Value{
		value.NewSymbol("$obj"): &value.MapValue{Value: map[value.Symbol]value.Value{
			value.NewSymbol("fox"):  &value.StringValue{Value: "FOX"},
			value.NewSymbol("fine"): &value.StringValue{Value: "FINE"},
		}},
		value.NewSymbol("$fun"): &value.BoolValue{Value: true},
	}}}

	{
		prefix, completions := inter.ReadEvalComplete("$obj f")
		assert.Equal(t, "$obj ", prefix)
		assert.Contains(t, completions, "fox")
		assert.Contains(t, completions, "fine")
	}

	{
		_, completions := inter.ReadEvalComplete("$obj $f")
		//assert.Equal(t, "$obj ", prefix)
		assert.Contains(t, completions, "$fun")
	}
}
