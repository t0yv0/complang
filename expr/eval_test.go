package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t0yv0/complang/value"
)

func TestEvalExpr(t *testing.T) {
	s := value.NewSymbol("foo")
	assert.Equal(t, &value.SymbolValue{Value: s}, EvalExpr(nil, &SymbolExpr{Symbol: s}))
}

func TestEvalQuery(t *testing.T) {
	// assume parsed "$obj foo"
	sq := &SymbolQuery{
		Expr: &RefExpr{
			Ref: value.NewSymbol("$obj"),
		},
		Symbol:       value.NewSymbol("foo"),
		SymbolOffset: 5,
	}
	foo1 := value.NewSymbol("foo1")
	foo2 := value.NewSymbol("foo2")
	env := map[value.Symbol]value.Value{
		value.NewSymbol("$obj"): &value.MapValue{
			Value: map[value.Symbol]value.Value{
				foo1: &value.StringValue{Value: "foo1value"},
				foo2: &value.StringValue{Value: "foo2value"},
			},
		},
	}
	res := EvalQuery(&value.MapEnv{env}, sq)
	assert.Equal(t, []value.Symbol{foo1, foo2}, res)
}
