package tools

import (
	"iter"
	"sync"
)

type LastSuccessIterator[T comparable] struct {
	lock  sync.RWMutex
	itens []T
}

func NewLastSuccessIterator[T comparable](itens ...T) *LastSuccessIterator[T] {
	return &LastSuccessIterator[T]{
		itens: itens,
	}
}

func (l *LastSuccessIterator[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		l.lock.RLock()
		lastIndex := 0
		for index, item := range l.itens {
			if !yield(item) {
				lastIndex = index
				break
			}
		}
		l.lock.RUnlock()
		if lastIndex == 0 {
			return
		}
		l.lock.Lock()
		var (
			before  = l.itens[:lastIndex]
			current = []T{l.itens[lastIndex]}
			after   = l.itens[lastIndex+1:]
		)
		l.itens = append(current, append(before, after...)...)
		l.lock.Unlock()
	}
}
