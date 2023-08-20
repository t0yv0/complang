package value

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"
)

var (
	symbolInternTable sync.Map
)

// [Value] space can be extended in Go by implementing this interface and using [CustomValue].
type ValueLike interface {
	// Implementations of Message should not have side-effects.
	Message(arg Value) Value
	// Programmable code completion, informing what can be passed to Message.
	CompleteSymbol(query Symbol) []Symbol
	// Run is a good place to perform side-effects if needed.
	Run() Value
	// The result of Show is how this value will be displayed to the user in the interpreter.
	Show() string
}

type Value interface {
	ValueLike
	valueMarker()
}

type Symbol interface {
	ValueLike
	symbolMarker()
	String() string
}

type sym struct {
	valueMarkerImpl
	noCompletions
	doesNotUnderstandAnything
	name string
}

var _ Symbol = (*sym)(nil)

func (*sym) symbolMarker()    {}
func (s *sym) String() string { return s.name }
func (s *sym) Run() Value     { return s }
func (s *sym) Show() string   { return s.name }

func NewSymbol(symbol string) Symbol {
	actual, _ := symbolInternTable.LoadOrStore(symbol, &sym{name: symbol})
	return actual.(Symbol)
}

type NullValue struct {
	valueMarkerImpl
	noCompletions
	doesNotUnderstandAnything
}

var _ Value = (*NullValue)(nil)

func (s *NullValue) Run() Value   { return s }
func (s *NullValue) Show() string { return "null" }

type BoolValue struct {
	valueMarkerImpl
	noCompletions
	doesNotUnderstandAnything
	Value bool
}

var _ Value = (*BoolValue)(nil)

func (s *BoolValue) Run() Value { return s }

func (s *BoolValue) Show() string {
	if s.Value {
		return "true"
	}
	return "false"
}

type ErrorValue struct {
	valueMarkerImpl
	noCompletions
	doesNotUnderstandAnything
	ErrorMessage string
}

var _ Value = (*ErrorValue)(nil)

func (s *ErrorValue) Run() Value { return s }

func (s *ErrorValue) Show() string { return fmt.Sprintf("ERROR: %s", s.ErrorMessage) }

type StringValue struct {
	valueMarkerImpl
	noCompletions
	doesNotUnderstandAnything
	Value string
}

var _ Value = (*StringValue)(nil)

func (s *StringValue) Run() Value { return s }

func (s *StringValue) Show() string { return fmt.Sprintf("%q", s.Value) }

type SymbolValue struct {
	valueMarkerImpl
	noCompletions
	doesNotUnderstandAnything
	Value Symbol
}

var _ Value = (*SymbolValue)(nil)

func (s *SymbolValue) Run() Value { return s }

func (s *SymbolValue) Show() string { return s.Value.String() }

type SliceValue struct {
	valueMarkerImpl
	noCompletions
	doesNotUnderstandAnything
	Value []Value
}

func (s *SliceValue) Run() Value { return s }

func (s *SliceValue) Show() string {
	if len(s.Value) == 0 {
		return "[]"
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "[")
	for i, v := range s.Value {
		if i > 0 {
			fmt.Fprintf(&buf, ", ")
		}
		fmt.Fprintf(&buf, v.Show())
	}
	fmt.Fprintf(&buf, "]")
	return buf.String()
}

type MapValue struct {
	valueMarkerImpl
	Value map[Symbol]Value
}

func (mv *MapValue) Message(arg Value) Value {
	if s, ok := arg.(*SymbolValue); ok {
		if v, ok := mv.Value[s.Value]; ok {
			return v
		}
	}
	return &ErrorValue{ErrorMessage: fmt.Sprintf("no map entries for key %s", arg.Show())}
}

func (mv *MapValue) CompleteSymbol(query Symbol) []Symbol {
	out := []Symbol{}
	for sym := range mv.Value {
		if strings.HasPrefix(sym.String(), query.String()) {
			out = append(out, sym)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].String() < out[j].String()
	})
	return out
}

func (mv *MapValue) Run() Value { return mv }

func (mv *MapValue) Show() string {
	out := []Symbol{}
	for sym := range mv.Value {
		out = append(out, sym)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].String() < out[j].String()
	})
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "{")
	for i, sym := range out {
		if i > 0 {
			fmt.Fprintf(&buf, ", ")
		}
		fmt.Fprintf(&buf, "%s: %s", sym.Show(), mv.Value[sym].Show())
	}
	fmt.Fprintf(&buf, "}")
	return buf.String()
}

var _ Value = (*MapValue)(nil)

type CustomValue struct {
	valueMarkerImpl
	ValueLike
}

var _ Value = (*CustomValue)(nil)

type valueMarkerImpl struct{}

func (*valueMarkerImpl) valueMarker() {}

type noCompletions struct{}

func (*noCompletions) CompleteSymbol(query Symbol) []Symbol {
	return nil
}

type doesNotUnderstandAnything struct{}

func (*doesNotUnderstandAnything) Message(value Value) Value {
	return &ErrorValue{ErrorMessage: fmt.Sprintf("object does not understand message %s", value.Show())}
}
