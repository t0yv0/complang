package complang

import (
	"context"
	"sync"
)

func LazyValue(f func() Value) Value {
	var once sync.Once
	var actual Value
	get := func() Value {
		once.Do(func() {
			actual = f()
		})
		return actual
	}
	return DeferredValue(get)
}

func DeferredValue(f func() Value) Value {
	return deferredValue(f)
}

type deferredValue func() Value

func (d deferredValue) Message(ctx context.Context, msg Value) Value {
	return d().Message(ctx, msg)
}
