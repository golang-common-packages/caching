package caching

import (
	"reflect"
	"sync"
)

// Item ...
type Item interface{}

// Queue ...
type Queue struct {
	items   sync.Map
	keys    []string
	maxSize int64
	curSize int64
	lock    sync.RWMutex
}

// NewQueue ...
func NewQueue(maxSize int64) *Queue {
	currentQueue := Queue{
		items:   sync.Map{},
		keys:    []string{},
		maxSize: maxSize,
		curSize: 0,
		lock:    sync.RWMutex{},
	}

	return &currentQueue
}

// Push ...
func (q *Queue) Push(k string, t Item) bool {
	itemSize := int64(reflect.Type.Size(t))
	if itemSize > q.maxSize || k == "" {
		return false
	}

	// Clean room for new items
	for q.curSize+itemSize > q.maxSize {
		q.Pop("")
	}

	q.lock.Lock()
	q.curSize = q.curSize + int64(reflect.Type.Size(t))
	q.keys = append(q.keys, k)
	q.items.LoadOrStore(k, t)
	q.lock.Unlock()
	return true
}

// Pop ...
func (q *Queue) Pop(k string) *Item {
	if q.IsEmpty() {
		return nil
	}

	var item interface{}
	var exits bool
	var idx int

	if k != "" {
		// Remove item by key
		item, exits = q.items.Load(k)
		for i, n := range q.keys {
			if k == n {
				idx = i
				break
			}
		}
	} else {
		// Remove first item
		key := q.keys[0]
		item, exits = q.items.Load(key)
		idx = 0
		if !exits {
			return nil
		}
	}

	q.lock.Lock()
	q.curSize = q.curSize - int64(reflect.Type.Size(item))
	q.items.Delete(k)
	q.keys = append(q.keys[:idx], q.keys[idx+1:]...)
	q.lock.Unlock()
	i, _ := (item).(Item)

	return &i
}

// Get ...
func (q *Queue) Get(k string) (Item, bool) {
	if q.IsEmpty() {
		return nil, false
	}

	item, exits := q.items.Load(k)
	if !exits {
		return nil, false
	}

	i, ok := (item).(Item)
	if !ok {
		return nil, false
	}

	return i, true
}

// Front ...
func (q *Queue) Front() *Item {
	if q.IsEmpty() {
		return nil
	}

	q.lock.RLock()
	key := q.keys[0]
	item, _ := q.items.Load(key)
	i, _ := item.(Item)
	q.lock.RUnlock()

	return &i
}

// Range ...
func (q *Queue) Range(fn func(key, value interface{}) bool) {
	q.items.Range(fn)
}

// IsEmpty ...
func (q *Queue) IsEmpty() bool {
	return len(q.keys) == 0
}

// Size ...
func (q *Queue) Size() int {
	return len(q.keys)
}
