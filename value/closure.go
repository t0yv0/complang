package value

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
		if v, ok := arg.(*SymbolValue); ok && v.Value == call {
			return c.Call(c.Env)
		}
		return &ErrorValue{ErrorMessage: "expecting call message"}
	}
	env := &extendedEnv{
		Env:    c.Env,
		symbol: c.Params[0],
		value:  arg,
	}
	if len(c.Params) == 1 {
		return c.Call(env)
	}
	return &CustomValue{ValueLike: &Closure{
		Env:    env,
		Params: c.Params[1:],
		Call:   c.Call,
	}}
}

func (c *Closure) CompleteSymbol(query Symbol) []Symbol { return nil }
func (c *Closure) Run() Value                           { return &CustomValue{ValueLike: c} }
func (c *Closure) Show() string                         { return "<closure>" }
