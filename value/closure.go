package value

import (
	"bytes"
	"fmt"
)

var (
	call = NewSymbol("call")
)

type Closure struct {
	Env    Env
	Params []Symbol
	Call   func(Env) Value
}

var _ ValueLike = (*Closure)(nil)

func (c *Closure) Message(arg Value) Value {
	if len(c.Params) == 0 {
		return &ErrorValue{ErrorMessage: "unexpected message call"}
	}
	env := &extendedEnv{
		Env:    c.Env,
		symbol: c.Params[0],
		value:  arg,
	}
	return &CustomValue{ValueLike: &Closure{
		Env:    env,
		Params: c.Params[1:],
		Call:   c.Call,
	}}
}

func (c *Closure) CompleteSymbol(query Symbol) []Symbol { return nil }

func (c *Closure) Run() Value {
	if len(c.Params) == 0 {
		return c.Call(c.Env)
	}
	return &CustomValue{ValueLike: c}
}

func (c *Closure) Show() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "<Closure")
	for i, p := range c.Params {
		if i > 0 {
			fmt.Fprintf(&buf, ",")
		} else {
			fmt.Fprintf(&buf, ":")
		}
		fmt.Fprintf(&buf, "%s", p.Show())
	}
	fmt.Fprintf(&buf, ">")
	return buf.String()
}
