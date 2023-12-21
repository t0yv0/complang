package expr

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	cl "github.com/t0yv0/complang"
)

func TestEvalExpr(t *testing.T) {
	ctx := context.Background()
	s := "foo"
	assert.Equal(t, cl.StringValue{Text: s},
		EvalExpr(ctx, nil, &SymbolExpr{Symbol: s}))
}

func TestEvalQuery(t *testing.T) {
	// assume parsed "$obj foo"
	sq := &SymbolQuery{
		Expr: &RefExpr{
			Ref: "$obj",
		},
		Symbol:       "foo",
		SymbolOffset: 5,
	}
	foo1 := "foo1"
	foo2 := "foo2"
	env := cl.NewMutableEnv()
	env.Bind("$obj", cl.MapValue(map[string]cl.Value{
		foo1: cl.StringValue{Text: "foo1value"},
		foo2: cl.StringValue{Text: "foo2value"},
	}))
	ctx := context.Background()
	result := []string{}
	EvalQuery(ctx, env, sq, func(_, s string) bool {
		result = append(result, s)
		return true
	})
	sort.Strings(result)
	assert.Equal(t, []string{foo1, foo2}, result)
}
