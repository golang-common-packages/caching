package caching

import (
	"reflect"
	"sync"
)

// Item the type of the queue
type Item interface{}

// ItemQueue the queue of Items
type ItemQueue struct {
	items   sync.Map
	keys    []string
	maxSize int64
	curSize int64
	lock    sync.RWMutex
}

// New creates a new ItemQueue
func ItemQueueNew(maxSize int64) *ItemQueue {
	s := ItemQueue{}
	s.items = sync.Map{}
	s.keys = []string{}
	s.maxSize = maxSize
	s.curSize = 0
	return &s
}

// Enqueue adds an Item to the end of the queue
func (s *ItemQueue) Push(k string, t Item) bool {
	itemSize := int64(reflect.Type.Size(t))
	if itemSize > s.maxSize || k == "" {
		return false
	}
	for s.curSize+itemSize > s.maxSize {
		s.Pop("")
	}
	s.lock.Lock()
	s.curSize = s.curSize + int64(reflect.Type.Size(t))
	s.keys = append(s.keys, k)
	s.items.LoadOrStore(k, t)
	s.lock.Unlock()
	return true
}

// Dequeue removes an Item from the start of the queue
func (s *ItemQueue) Pop(k string) *Item {
	if s.IsEmpty() {
		return nil
	}
	var item interface{}
	var exits bool
	var idx int
	if k != "" {
		item, exits = s.items.Load(k)
		for i, n := range s.keys {
			if k == n {
				idx = i
				break
			}
		}
	} else {
		key := s.keys[0]
		item, exits = s.items.Load(key)
		idx = 0
	}
	if exits {
		s.lock.Lock()
		s.curSize = s.curSize - int64(reflect.Type.Size(item))
		s.items.Delete(k)
		s.keys = append(s.keys[:idx], s.keys[idx+1:]...)
		s.lock.Unlock()
		i, _ := (item).(Item)
		return &i
	} else {
		return nil
	}
}

func (s *ItemQueue) Get(k string) (Item, bool) {
	if s.IsEmpty() {
		return nil, false
	}
	item, exits := s.items.Load(k)
	if exits {
		i, ok := (item).(Item)
		if !ok {
			return nil, false
		}
		return i, true
	} else {
		return nil, false
	}
}

// Front returns the item next in the queue, without removing it
func (s *ItemQueue) Front() *Item {
	if s.IsEmpty() {
		return nil
	}
	s.lock.RLock()
	key := s.keys[0]
	item, _ := s.items.Load(key)
	i, _ := item.(Item)
	s.lock.RUnlock()
	return &i
}

func (s *ItemQueue) Range(fn func(key, value interface{}) bool) {
	s.items.Range(fn)
}

// IsEmpty returns true if the queue is empty
func (s *ItemQueue) IsEmpty() bool {
	return len(s.keys) == 0
}

// Size returns the number of Items in the queue
func (s *ItemQueue) Size() int {
	return len(s.keys)
}
