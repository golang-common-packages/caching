package caching

import (
	"errors"
	"reflect"
	"sync"
)

// Queue ...
type Queue struct {
	items            sync.Map
	keys             []string
	queueSize        int64
	queueCurrentSize int64
	rwMutex          sync.RWMutex
}

// NewQueue ...
func NewQueue(queueSize int64) *Queue {
	currentQueue := Queue{
		items:            sync.Map{},
		keys:             []string{},
		queueSize:        queueSize,
		queueCurrentSize: 0,
		rwMutex:          sync.RWMutex{},
	}

	return &currentQueue
}

// Enqueue ...
func (q *Queue) Enqueue(key string, value interface{}) error {
	itemSize := int64(reflect.Type.Size(value))
	if itemSize > q.queueSize || key == "" {
		return errors.New("key is empty or queue not enough space")
	}

	// Clean space for new item
	for q.queueCurrentSize+itemSize > q.queueSize {
		q.Dequeue("")
	}

	q.rwMutex.Lock()
	q.queueCurrentSize = q.queueCurrentSize + int64(reflect.Type.Size(value))
	q.keys = append(q.keys, key)
	q.items.LoadOrStore(key, value)
	q.rwMutex.Unlock()

	return nil
}

// Dequeue ...
func (q *Queue) Dequeue(key string) error {
	if q.IsEmpty() {
		return errors.New("the queue is empty")
	}

	var (
		item     interface{}
		exits    bool
		keyIndex int
	)

	if key != "" {
		// Load item by key
		item, exits = q.items.Load(key)
		for i, k := range q.keys {
			if key == k {
				keyIndex = i
				break
			}
		}
	} else {
		// Load the first item
		item, exits = q.items.Load(q.keys[0])
		keyIndex = 0
	}

	if !exits {
		return errors.New("this key does not exits")
	}

	q.rwMutex.Lock()
	q.queueCurrentSize = q.queueCurrentSize - int64(reflect.Type.Size(item))
	q.items.Delete(key)
	q.keys = append(q.keys[:keyIndex], q.keys[keyIndex+1:]...) //Update keys slice after remove this key from items map
	q.rwMutex.Unlock()

	return nil
}

// Get ...
func (q *Queue) Get() (interface{}, error) {
	if q.IsEmpty() {
		return nil, errors.New("queue is empty")
	}

	q.rwMutex.RLock()
	key := q.keys[0]
	item, exits := q.items.Load(key)
	q.rwMutex.RUnlock()
	if !exits {
		return nil, errors.New("this key does not exits")
	}

	return item, nil
}

// GetByKey ...
func (q *Queue) GetByKey(k string) (interface{}, error) {
	if q.IsEmpty() {
		return nil, errors.New("queue is empty")
	}

	q.rwMutex.RLock()
	item, exits := q.items.Load(k)
	q.rwMutex.RUnlock()
	if !exits {
		return nil, errors.New("this key does not exits")
	}

	return item, nil
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

// Capacity ...
func (q *Queue) Capacity() int64 {
	return q.queueCurrentSize
}
