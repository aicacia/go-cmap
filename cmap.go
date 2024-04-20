package cmap

import (
	"sync"
	"sync/atomic"
)

type CMap[K, V any] struct {
	sync.Map
	count atomic.Int64
}

func New[K, V any]() CMap[K, V] {
	return CMap[K, V]{}
}

func (m *CMap[K, V]) SetIfAbsent(key K, value V) bool {
	if _, isOld := m.Map.LoadOrStore(key, value); !isOld {
		m.count.Add(1)
		return true
	} else {
		return false
	}
}

func (m *CMap[K, V]) Set(key K, value V) bool {
	if _, isOld := m.Map.Swap(key, value); !isOld {
		m.count.Add(1)
		return true
	} else {
		return false
	}
}

func (m *CMap[K, V]) Has(key K) bool {
	_, ok := m.Map.Load(key)
	return ok
}

func (m *CMap[K, V]) IsEmpty() bool {
	return m.Count() == 0
}

func (m *CMap[K, V]) Get(key K) (V, bool) {
	if value, ok := m.Map.Load(key); ok {
		return value.(V), true
	} else {
		return *new(V), false
	}
}

func (m *CMap[K, V]) GetOrSet(key K, value V) V {
	result, isOld := m.Map.LoadOrStore(key, value)
	if !isOld {
		m.count.Add(1)
	}
	return result.(V)
}

func (m *CMap[K, V]) Delete(key K) bool {
	if _, isOld := m.Map.LoadAndDelete(key); isOld {
		m.count.Add(-1)
		return true
	} else {
		return false
	}
}

func (m *CMap[K, V]) Remove(key K) bool {
	return m.Delete(key)
}

type Entry[K, V any] struct {
	Key K
	Val V
}

func (m *CMap[K, V]) Iter() chan Entry[K, V] {
	ch := make(chan Entry[K, V])
	go func() {
		m.Map.Range(func(key, value any) bool {
			ch <- Entry[K, V]{
				Key: key.(K),
				Val: value.(V),
			}
			return true
		})
		close(ch)
	}()
	return ch
}

func (m *CMap[K, V]) Keys() chan K {
	ch := make(chan K)
	go func() {
		m.Map.Range(func(key, _ any) bool {
			ch <- key.(K)
			return true
		})
		close(ch)
	}()
	return ch
}

func (m *CMap[K, V]) Values() chan V {
	ch := make(chan V)
	go func() {
		m.Map.Range(func(_, value any) bool {
			ch <- value.(V)
			return true
		})
		close(ch)
	}()
	return ch
}

func (m *CMap[K, V]) Count() int64 {
	return m.count.Load()
}

func (m *CMap[K, V]) Len() int {
	return int(m.Count())
}

func (m *CMap[K, V]) Clear() {
	m.Map.Range(func(key, _ any) bool {
		m.Delete(key.(K))
		return true
	})
}
