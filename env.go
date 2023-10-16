package complang

import (
	"sort"
)

type Env interface {
	Symbols() []string
	Lookup(symbol string) (Value, bool)
}

type MutableEnv interface {
	Env
	Bind(symbol string, value Value)
	Unbind(symbol string)
}

func NewMutableEnv() MutableEnv {
	return mutableEnv{make(mapEnv)}
}

type mutableEnv struct {
	mapEnv
}

var _ MutableEnv = (*mutableEnv)(nil)

func (me mutableEnv) Bind(symbol string, value Value) {
	if me.mapEnv == nil {
		me.mapEnv = map[string]Value{}
	}
	me.mapEnv[symbol] = value
}

func (me mutableEnv) Unbind(symbol string) {
	if me.mapEnv != nil {
		delete(me.mapEnv, symbol)
	}
}

type mapEnv map[string]Value

func (me mapEnv) Symbols() []string {
	var ss []string
	for x := range me {
		ss = append(ss, x)
	}
	sort.Strings(ss)
	return ss
}

func (me mapEnv) Lookup(s string) (Value, bool) {
	v, ok := me[s]
	return v, ok
}

type extendedEnv struct {
	Env
	symbol string
	value  Value
}

func (ee *extendedEnv) Symbols() []string {
	ss := ee.Env.Symbols()
	for _, s := range ss {
		if s == ee.symbol {
			return ss
		}
	}
	return append(ss, ee.symbol)
}

func (ee *extendedEnv) Lookup(s string) (Value, bool) {
	if s == ee.symbol {
		return ee.value, true
	}
	return ee.Env.Lookup(s)
}
