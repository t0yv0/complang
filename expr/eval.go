package expr

import (
	"fmt"

	"github.com/t0yv0/complang/value"
)

func EvalExpr(env map[value.Symbol]value.Value, expr Expr) value.Value {
	switch expr := expr.(type) {
	case *NullExpr:
		return &value.NullValue{}
	case *BoolExpr:
		return &value.BoolValue{Value: expr.Bool}
	case *SymbolExpr:
		return &value.SymbolValue{Value: expr.Symbol}
	case *StringExpr:
		return &value.StringValue{Value: expr.String}
	case *RefExpr:
		v, ok := env[expr.Ref]
		if ok {
			return v
		}
		return &value.ErrorValue{ErrorMessage: fmt.Sprintf("unbound symbol: %s", expr.Ref)}
	case *MessageExpr:
		receiver := EvalExpr(env, expr.Receiver)
		message := EvalExpr(env, expr.Message)
		return receiver.Message(message)
	case *LambdaBlockExpr:
		return &value.CustomValue{ValueLike: &Closure{
			params: expr.Symbols,
			body:   expr.Body,
			env:    env,
		}}
	default:
		panic("EvalExpr is incomplete")
	}
}

type Closure struct {
	env    map[value.Symbol]value.Value
	params []value.Symbol
	body   Expr
}

var (
	call = value.NewSymbol("call")
)

func (c *Closure) Message(arg value.Value) value.Value {
	if len(c.params) == 0 {
		if v, ok := arg.(*value.SymbolValue); ok && v.Value == call {
			return EvalExpr(c.env, c.body)
		}
		return &value.ErrorValue{ErrorMessage: "expecting call message"}
	}
	env := map[value.Symbol]value.Value{}
	for k, v := range c.env {
		env[k] = v
	}
	env[c.params[0]] = arg
	if len(c.params) == 1 {
		return EvalExpr(env, c.body)
	}
	return &value.CustomValue{ValueLike: &Closure{
		env:    env,
		params: c.params[1:],
		body:   c.body,
	}}
}

func (c *Closure) CompleteSymbol(query value.Symbol) []value.Symbol { return nil }
func (c *Closure) Run() value.Value                                 { return &value.CustomValue{ValueLike: c} }
func (c *Closure) Show() string                                     { return "<closure>" }

var _ value.ValueLike = (*Closure)(nil)

func EvalStmt(env map[value.Symbol]value.Value, stmt Stmt) {
	switch stmt := stmt.(type) {
	case *ExprStmt:
		v := EvalExpr(env, stmt.Expr)
		v = v.Run() // run side-effects
		fmt.Println(v.Show())
	case *AssignStmt:
		v := EvalExpr(env, stmt.Expr)
		v = v.Run() // run side-effects
		env[stmt.Ref] = v
	default:
		panic("EvalStmt is incomplete")
	}
}

func EvalQuery(env map[value.Symbol]value.Value, q Query) []value.Symbol {
	switch q := q.(type) {
	case *SymbolQuery:
		v := EvalExpr(env, q.Expr)
		return v.CompleteSymbol(q.Symbol)
	case *RefQuery:
		v := &value.MapValue{Value: env}
		return v.CompleteSymbol(q.Ref)
	default:
		panic(fmt.Sprintf("EvalQuery is incomplete, got %#T", q))
	}
}
