package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t0yv0/complang/expr"
)

func TestParseStmt(t *testing.T) {
	s, err := ParseStmt(`$x = "$foo"`)
	assert.NoError(t, err)
	t.Logf("#%T", s)
	s, ok := s.(*expr.AssignStmt)
	assert.True(t, ok)
	t.Logf("%v", s)
}
