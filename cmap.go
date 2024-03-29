package cmap

import (
	"sync"
	"sync/atomic"
)

type CMap[T any] struct {
	sync.Map
	count atomic.Int64
}

func New[T any]() CMap[T] {
	return CMap[T]{}
}

func (m *CMap[T]) SetIfAbsent(key string, value T) bool {
	if _, isOld := m.Map.LoadOrStore(key, value); !isOld {
		m.count.Add(1)
		return true
	} else {
		return false
	}
}

func (m *CMap[T]) Set(key string, value T) bool {
	if _, isOld := m.Map.Swap(key, value); !isOld {
		m.count.Add(1)
		return true
	} else {
		return false
	}
}

func (m *CMap[T]) Has(key string) bool {
	_, ok := m.Map.Load(key)
	return ok
}

func (m *CMap[T]) IsEmpty() bool {
	return m.Count() == 0
}

func (m *CMap[T]) Get(key string) (T, bool) {
	if value, ok := m.Map.Load(key); ok {
		return value.(T), true
	} else {
		return *new(T), false
	}
}

func (m *CMap[T]) GetOrSet(key string, value T) T {
	result, isOld := m.Map.LoadOrStore(key, value)
	if !isOld {
		m.count.Add(1)
	}
	return result.(T)
}

func (m *CMap[T]) Delete(key string) bool {
	if _, isOld := m.Map.LoadAndDelete(key); isOld {
		m.count.Add(-1)
		return true
	} else {
		return false
	}
}

func (m *CMap[T]) Remove(key string) bool {
	return m.Delete(key)
}

type Entry[T any] struct {
	Key string
	Val T
}

func (m *CMap[T]) Iter() chan Entry[T] {
	ch := make(chan Entry[T])
	go func() {
		m.Map.Range(func(key, value any) bool {
			ch <- Entry[T]{
				Key: key.(string),
				Val: value.(T),
			}
			return true
		})
		close(ch)
	}()
	return ch
}

func (m *CMap[T]) Keys() chan string {
	ch := make(chan string)
	go func() {
		m.Map.Range(func(key, _ any) bool {
			ch <- key.(string)
			return true
		})
		close(ch)
	}()
	return ch
}

func (m *CMap[T]) Values() chan T {
	ch := make(chan T)
	go func() {
		m.Map.Range(func(_, value any) bool {
			ch <- value.(T)
			return true
		})
		close(ch)
	}()
	return ch
}

func (m *CMap[T]) Count() int64 {
	return m.count.Load()
}

func (m *CMap[T]) Len() int {
	return int(m.Count())
}

func (m *CMap[T]) Clear() {
	m.Map.Range(func(key, _ any) bool {
		m.Delete(key.(string))
		return true
	})
}
