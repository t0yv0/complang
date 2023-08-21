package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t0yv0/complang/expr"
)

func TestParseStmt(t *testing.T) {
	s, err := ParseStmt(`$x = "$foo"`)
	assert.NoError(t, err)
	stmt, ok := s.(*expr.AssignStmt)
	assert.True(t, ok)
	assert.Equal(t, "$x", stmt.Ref.String())
	str, ok := stmt.Expr.(*expr.StringExpr)
	assert.True(t, ok)
	assert.Equal(t, "$foo", str.String)
}

func TestParseQuery(t *testing.T) {
	t.Run("SymbolQuery", func(t *testing.T) {
		for _, code := range []string{"$obj f", "$v = $obj f"} {
			q, err := ParseQuery(code)
			assert.NoError(t, err)
			sq, ok := q.(*expr.SymbolQuery)
			assert.True(t, ok)
			re, ok := sq.Expr.(*expr.RefExpr)
			assert.True(t, ok)
			assert.Equal(t, "$obj", re.Ref.String())
			assert.Equal(t, "f", sq.Symbol.String())
		}
	})
	t.Run("RefQuery", func(t *testing.T) {
		{
			q, err := ParseQuery("$f")
			assert.NoError(t, err)
			rq, ok := q.(*expr.RefQuery)
			assert.True(t, ok)
			assert.Equal(t, "$f", rq.Ref.String())
		}
		{
			q, err := ParseQuery("$obj $f")
			assert.NoError(t, err)
			rq, ok := q.(*expr.RefQuery)
			assert.True(t, ok)
			assert.Equal(t, "$f", rq.Ref.String())
		}
	})
}
