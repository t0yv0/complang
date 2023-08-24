package value

import "sort"

type Env interface {
	Symbols() []Symbol
	Lookup(Symbol) (Value, bool)
}

type MapEnv struct {
	SymbolMap map[Symbol]Value
}

func (me *MapEnv) Symbols() []Symbol {
	var ss []Symbol
	for x := range me.SymbolMap {
		ss = append(ss, x)
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].String() < ss[j].String()
	})
	return ss
}

func (me *MapEnv) Lookup(s Symbol) (Value, bool) {
	v, ok := me.SymbolMap[s]
	return v, ok
}

type extendedEnv struct {
	Env
	symbol Symbol
	value  Value
}

func (ee *extendedEnv) Lookup(s Symbol) (Value, bool) {
	if s == ee.symbol {
		return ee.value, true
	}
	return ee.Env.Lookup(s)
}
