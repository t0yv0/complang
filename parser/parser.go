package parser

import (
	"fmt"
	"strings"

	"github.com/t0yv0/complang/expr"
	"github.com/t0yv0/complang/value"
)

func ParseExpr(code string) (expr.Expr, error) {
	tokens, err := tokenize(code)
	if err != nil {
		return nil, err
	}
	e, rest := parseExpr(tokens)
	if e == nil && len(rest) > 0 {
		return nil, fmt.Errorf("could not parse expression")
	}
	return e, nil
}

func ParseStmt(code string) (expr.Stmt, error) {
	tokens, err := tokenize(code)
	if err != nil {
		return nil, err
	}
	e, rest := parseStmt(tokens)
	if e == nil && len(rest) > 0 {
		return nil, fmt.Errorf("could not parse expression")
	}
	return e, nil
}

func ParseQuery(code string) (expr.Query, error) {
	tokens, err := tokenize(code)
	if err != nil {
		return nil, err
	}
	e, rest := parseQuery(tokens)
	if e == nil || len(rest) > 0 {
		return nil, fmt.Errorf("could not parse expression in query")
	}
	return e, nil
}

func parseExpr(tokens []token) (expr.Expr, []token) {
	e, tokens := parseSimpleExpr(tokens)
	if e == nil {
		return nil, tokens
	}
	for {
		subE, rest := parseSimpleExpr(tokens)
		tokens = rest
		if subE == nil {
			return e, tokens
		} else {
			e = &expr.MessageExpr{Receiver: e, Message: subE}
		}
	}
}

func parseSimpleExpr(tokens []token) (expr.Expr, []token) {
	if len(tokens) == 0 {
		return nil, tokens
	}
	if tokens[0].t == byte('(') {
		e, rest := parseExpr(tokens[1:])
		if len(rest) > 0 && rest[0].t == byte(')') {
			rest = rest[1:]
		}
		return e, rest
	}
	switch t := tokens[0].t.(type) {
	case value.Symbol:
		if isRef(t) {
			return &expr.RefExpr{Ref: t}, tokens[1:]
		}
		return &expr.SymbolExpr{
			Symbol: t,
			Offset: tokens[0].offset,
		}, tokens[1:]
	case string:
		return &expr.StringExpr{String: t}, tokens[1:]
	case bool:
		return &expr.BoolExpr{Bool: t}, tokens[1:]
	case nil:
		return &expr.NullExpr{}, tokens[1:]
	default:
		return nil, tokens
	}
}

func isRef(s value.Symbol) bool {
	return strings.HasPrefix(s.Show(), "$")
}

func parseStmt(tokens []token) (expr.Stmt, []token) {
	if len(tokens) == 0 {
		return nil, tokens
	}
	if s, ok := tokens[0].t.(value.Symbol); ok && isRef(s) {
		if len(tokens) > 1 && tokens[1].t == byte('=') {
			e, rest := parseExpr(tokens[2:])
			if e != nil {
				return &expr.AssignStmt{
					Ref:  s,
					Expr: e,
				}, rest
			}
		}
	}
	e, rest := parseExpr(tokens)
	if e != nil {
		return &expr.ExprStmt{
			Expr: e,
		}, rest
	}
	return nil, tokens
}

func parseQuery(tokens []token) (expr.Query, []token) {
	e, rest := parseExpr(tokens)
	if e == nil {
		return nil, tokens
	}
	switch e := e.(type) {
	case *expr.MessageExpr:
		switch s := e.Message.(type) {
		case *expr.SymbolExpr:
			return &expr.SymbolQuery{
				Expr:         e.Receiver,
				Symbol:       s.Symbol,
				SymbolOffset: s.Offset,
			}, rest
		}
	}
	return nil, tokens
}
