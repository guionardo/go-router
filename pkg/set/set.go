package set

import (
	"iter"
	"sync"
)

type Set[T comparable] struct {
	set  map[T]struct{}
	lock sync.RWMutex
}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		set: make(map[T]struct{}),
	}
}

func (s *Set[T]) Set(value T) {
	s.lock.Lock()
	s.set[value] = struct{}{}
	s.lock.Unlock()
}

func (s *Set[T]) Unset(value T) {
	s.lock.Lock()
	delete(s.set, value)
	s.lock.Unlock()
}

func (s *Set[T]) IsSet(value T) bool {
	s.lock.RLock()
	_, ok := s.set[value]
	s.lock.RUnlock()
	return ok
}

func (s *Set[T]) Values() (values []T) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k := range s.set {
		values = append(values, k)
	}
	return values
}

func (s *Set[T]) Range() iter.Seq[T] {
	return func(yield func(T) bool) {
		s.lock.RLock()
		defer s.lock.RUnlock()
		for k := range s.set {
			if !yield(k) {
				return
			}
		}
	}
}
