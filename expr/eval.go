package expr

import (
	"context"
	"fmt"

	cl "github.com/t0yv0/complang"
)

func EvalExpr(ctx context.Context, env cl.Env, expr Expr) cl.Value {
	switch expr := expr.(type) {
	case *NullExpr:
		return cl.NullValue{}
	case *BoolExpr:
		return cl.BoolValue{Bool: expr.Bool}
	case *StringExpr:
		return cl.StringValue{Text: expr.String}
	case *NumExpr:
		return cl.NumValue{Num: expr.Number}
	case *SymbolExpr:
		return cl.StringValue{Text: expr.Symbol}
	case *RefExpr:
		v, ok := env.Lookup(expr.Ref)
		if ok {
			return v
		}
		return cl.Error{ErrorMessage: fmt.Sprintf("unbound symbol: %s", expr.Ref)}
	case *MessageExpr:
		receiver := EvalExpr(ctx, env, expr.Receiver)
		message := EvalExpr(ctx, env, expr.Message)
		return receiver.Message(ctx, message)
	case *LambdaBlockExpr:
		body := expr.Body
		return cl.Closure{
			Env:    env,
			Params: expr.Symbols,
			Call: func(ctx context.Context, env cl.Env) cl.Value {
				return EvalExpr(ctx, env, body)
			},
		}
	default:
		panic("EvalExpr is incomplete")
	}
}

func EvalStmt(ctx context.Context, env cl.MutableEnv, stmt Stmt) {
	switch stmt := stmt.(type) {
	case *ExprStmt:
		v := EvalExpr(ctx, env, stmt.Expr)
		v = cl.Run(ctx, v) // run side-effects
		fmt.Println(cl.Show(ctx, v))
	case *AssignStmt:
		v := EvalExpr(ctx, env, stmt.Expr)
		v = cl.Run(ctx, v) // run side-effects
		env.Bind(stmt.Ref, v)
	default:
		panic("EvalStmt is incomplete")
	}
}

func EvalQuery(ctx context.Context, env cl.Env, q Query, complete func(string, string) bool) {
	switch q := q.(type) {
	case *SymbolQuery:
		v := EvalExpr(ctx, env, q.Expr)
		cl.Complete(ctx, v, cl.CompleteRequest{
			Query:    q.Symbol,
			Receiver: complete,
		})
	case *RefQuery:
		cl.Complete(ctx, envToMap(env), cl.CompleteRequest{
			Query:    q.Ref,
			Receiver: complete,
		})
	default:
		panic(fmt.Sprintf("EvalQuery is incomplete, got %#T", q))
	}
}

func envToMap(env cl.Env) cl.Value {
	v := make(cl.MapValue)
	for _, s := range env.Symbols() {
		if sv, ok := env.Lookup(s); ok {
			v[s] = sv
		}
	}
	return v
}
