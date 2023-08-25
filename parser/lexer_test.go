package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/t0yv0/complang/value"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	type testCase struct {
		s      string
		tokens []any
		err    error
	}

	testCases := []testCase{
		{
			s:      "null",
			tokens: []any{nil},
		},
		{
			s:      "true",
			tokens: []any{true},
		},
		{
			s:      "false",
			tokens: []any{false},
		},
		{
			s:      "  false",
			tokens: []any{false},
		},
		{
			s:      "true null false",
			tokens: []any{true, nil, false},
		},
		{
			s:      "sym",
			tokens: []any{value.NewSymbol("sym")},
		},
		{
			s:      "$ref",
			tokens: []any{value.NewSymbol("$ref")},
		},
		{
			s:      `"foo"`,
			tokens: []any{"foo"},
		},
		{
			s:      `"foo""" "b a r\n"`,
			tokens: []any{"foo", "", "b a r\n"},
		},
		{
			s: `$ref = ("foo")`,
			tokens: []any{
				value.NewSymbol("$ref"),
				byte('='),
				byte('('),
				"foo",
				byte(')'),
			},
		},
		{
			s:      "123",
			tokens: []any{123},
			err:    fmt.Errorf("unexpected '1'"),
		},
		{
			s:      "foo:bar/baz",
			tokens: []any{value.NewSymbol("foo:bar/baz")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			actual, err := tokenize(tc.s)
			if err != nil {
				if tc.err != nil {
					if err.Error() != tc.err.Error() {
						t.Error(err)
						t.FailNow()
					}
				} else {
					t.Error(err)
					t.FailNow()
				}
			} else {
				actualValues := []any{}
				for _, v := range actual {
					actualValues = append(actualValues, v.t)
				}
				if !reflect.DeepEqual(actualValues, tc.tokens) {
					t.Errorf("unexpected %v, expecting %v", actual, tc.tokens)
				}
			}
		})
	}
}

func TestOffsets(t *testing.T) {
	source := `$obj fld "string\n" subf`
	tokens, err := tokenize(source)
	assert.NoError(t, err)
	assert.Equal(t, value.NewSymbol("subf"), tokens[len(tokens)-1].t)
	assert.Equal(t, "subf", source[tokens[len(tokens)-1].offset:])
}
