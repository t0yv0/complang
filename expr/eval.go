package expr

import (
	"fmt"

	"github.com/t0yv0/complang/value"
)

func EvalExpr(env value.Env, expr Expr) value.Value {
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
		v, ok := env.Lookup(expr.Ref)
		if ok {
			return v
		}
		return &value.ErrorValue{ErrorMessage: fmt.Sprintf("unbound symbol: %s", expr.Ref)}
	case *MessageExpr:
		receiver := EvalExpr(env, expr.Receiver)
		message := EvalExpr(env, expr.Message)
		return receiver.Message(message)
	case *LambdaBlockExpr:
		body := expr.Body
		closure := &value.Closure{
			Env:    env,
			Params: expr.Symbols,
			Call: func(env value.Env) value.Value {
				return EvalExpr(env, body)
			},
		}
		return &value.CustomValue{ValueLike: closure}
	default:
		panic("EvalExpr is incomplete")
	}
}

func EvalStmt(env *value.MapEnv, stmt Stmt) {
	switch stmt := stmt.(type) {
	case *ExprStmt:
		v := EvalExpr(env, stmt.Expr)
		v = v.Run() // run side-effects
		fmt.Println(v.Show())
	case *AssignStmt:
		v := EvalExpr(env, stmt.Expr)
		v = v.Run() // run side-effects
		env.SymbolMap[stmt.Ref] = v
	default:
		panic("EvalStmt is incomplete")
	}
}

func EvalQuery(env value.Env, q Query) []value.Symbol {
	switch q := q.(type) {
	case *SymbolQuery:
		v := EvalExpr(env, q.Expr)
		return v.CompleteSymbol(q.Symbol)
	case *RefQuery:
		v := &value.MapValue{Value: map[value.Symbol]value.Value{}}
		for _, s := range env.Symbols() {
			if sv, ok := env.Lookup(s); ok {
				v.Value[s] = sv
			}
		}
		return v.CompleteSymbol(q.Ref)
	default:
		panic(fmt.Sprintf("EvalQuery is incomplete, got %#T", q))
	}
}
