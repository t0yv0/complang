package complang

import (
	"bytes"
	"context"
	"fmt"
)

type Value interface {
	Message(context.Context, Value) Value
}

// :show -- objects respond to :show with a string value to customize how they will be displayed.
type ShowMessage struct{}

func (ShowMessage) Message(ctx context.Context, v Value) Value {
	switch v.(type) {
	case ShowMessage:
		return StringValue{":show"}
	default:
		return DoNotUnderstandError(ctx, v)
	}
}

// :run -- objects are only allowed to do side-effects when responding to the :run message.
type RunMessage struct{}

func (RunMessage) Message(ctx context.Context, v Value) Value {
	switch v.(type) {
	case ShowMessage:
		return StringValue{":run"}
	default:
		return DoNotUnderstandError(ctx, v)
	}
}

// :comlete -- objects respond to :complete to customize code completion.
type CompleteRequest struct {
	Query    string
	Receiver func(query string, match string) bool
}

func (CompleteRequest) Message(ctx context.Context, v Value) Value {
	switch v.(type) {
	case ShowMessage:
		return StringValue{":complete"}
	default:
		return DoNotUnderstandError(ctx, v)
	}
}

type Error struct {
	ErrorMessage string
}

func (x Error) Message(_ context.Context, v Value) Value {
	switch v.(type) {
	case ShowMessage:
		return StringValue{fmt.Sprintf("ERROR: %s", x.ErrorMessage)}
	default:
		// Errors are self-evaluating to float out of expressions.
		return x
	}
}

func DoNotUnderstandError(ctx context.Context, message Value) Value {
	return &Error{fmt.Sprintf("object does not understand %s", Show(ctx, message))}
}

type StringValue struct {
	Text string
}

func (x StringValue) Message(ctx context.Context, v Value) Value {
	switch v.(type) {
	case ShowMessage:
		return StringValue{fmt.Sprintf("%q", x.Text)}
	default:
		return DoNotUnderstandError(ctx, v)
	}
}

func Show(ctx context.Context, v Value) string {
	switch x := v.Message(ctx, ShowMessage{}).(type) {
	case StringValue:
		return x.Text
	default:
		return Show(ctx, Error{"object does not respond to :show properly"})
	}
}

func Run(ctx context.Context, v Value) Value {
	return v.Message(ctx, RunMessage{})
}

func Complete(ctx context.Context, v Value, req CompleteRequest) {
	v.Message(ctx, req)
}

type NullValue struct{}

func (x NullValue) Message(ctx context.Context, v Value) Value {
	switch v.(type) {
	case ShowMessage:
		return StringValue{"null"}
	default:
		return DoNotUnderstandError(ctx, v)
	}
}

type BoolValue struct {
	Bool bool
}

func (x BoolValue) Message(ctx context.Context, v Value) Value {
	switch v.(type) {
	case ShowMessage:
		if x.Bool {
			return StringValue{"true"}
		}
		return StringValue{"false"}
	default:
		return DoNotUnderstandError(ctx, v)
	}
}

type SliceValue []Value

func (x SliceValue) Message(ctx context.Context, v Value) Value {
	switch v := v.(type) {
	case ShowMessage:
		if len(x) == 0 {
			return StringValue{"[]"}
		}
		return StringValue{fmt.Sprintf("[len=%d]", len(x))}
	default:
		return DoNotUnderstandError(ctx, v)
	}
}

type MapValue map[string]Value

func (x MapValue) Message(ctx context.Context, v Value) Value {
	switch v := v.(type) {
	case ShowMessage:
		if len(x) == 0 {
			return StringValue{"{}"}
		}
		return StringValue{fmt.Sprintf("{len=%d}", len(x))}
	case CompleteRequest:
		for k := range x {
			if !v.Receiver(v.Query, k) {
				break
			}
		}
		return NullValue{}
	default:
		return DoNotUnderstandError(ctx, v)
	}
}

type Closure struct {
	Env    Env
	Params []string
	Call   func(context.Context, Env) Value
}

// func (c *Closure) Message(arg Value) Value {
// 	if len(c.Params) == 0 {
// 		return &ErrorValue{ErrorMessage: "unexpected message call"}
// 	}
// 	env := &extendedEnv{
// 		Env:    c.Env,
// 		symbol: c.Params[0],
// 		value:  arg,
// 	}
// 	return &CustomValue{ValueLike: &Closure{
// 		Env:    env,
// 		Params: c.Params[1:],
// 		Call:   c.Call,
// 	}}
// }

func (c Closure) Message(ctx context.Context, msg Value) Value {
	switch msg := msg.(type) {
	case ShowMessage:
		return StringValue{c.show()}
	case RunMessage:
		return c.run(ctx)
	case Error:
		return msg
	default:
		if len(c.Params) == 0 {
			return DoNotUnderstandError(ctx, msg)
		}
		return Closure{
			Env: &extendedEnv{
				Env:    c.Env,
				symbol: c.Params[0],
				value:  msg,
			},
			Params: c.Params[1:],
			Call:   c.Call,
		}
	}
}

func (c Closure) run(ctx context.Context) Value {
	if len(c.Params) == 0 {
		return c.Call(ctx, c.Env)
	}
	return c
}

func (c Closure) show() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "<Closure")
	for i, p := range c.Params {
		if i > 0 {
			fmt.Fprintf(&buf, ",")
		} else {
			fmt.Fprintf(&buf, ":")
		}
		fmt.Fprintf(&buf, "%s", p)
	}
	fmt.Fprintf(&buf, ">")
	return buf.String()
}
