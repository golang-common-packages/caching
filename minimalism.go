package caching

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
	"github.com/golang-common-packages/queue"
)

// MinimalismClient ...
type MinimalismClient struct {
	items *queue.Queue
	close chan struct{}
}

// NewMinimalism ...
func NewMinimalism(config *Config) ICaching {
	currentSession := &MinimalismClient{queue.New(config.Minimalism.CacheSize), make(chan struct{})}

	go func() {
		ticker := time.NewTicker(config.Minimalism.CleaningInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				now := time.Now().UnixNano()

				currentSession.items.Range(func(key, value interface{}) bool {
					item := value.(MinimalismItem)

					if item.expires < now {
						k, _ := key.(string)
						currentSession.items.Dequeue(k)
					}

					return true
				})

			case <-currentSession.close:
				return
			}
		}
	}()

	return currentSession
}

// Middleware ...
func (ml *MinimalismClient) Middleware(hash hash.IHash) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			key := hash.SHA512(token)

			if val, err := ml.Get(key); err != nil {
				log.Printf("Can not get accesstoken from redis in redis middleware: %s", err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			} else if val == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// GetByKey ...
func (ml *MinimalismClient) Get(key string) (interface{}, error) {
	obj, err := ml.items.GetByKey(key)
	if err != nil {
		return nil, errors.New("item with that key does not exist")
	}

	item, ok := obj.(MinimalismItem)
	if !ok {
		return nil, errors.New("can not map object to MinimalismItem model")
	}

	if item.expires > 0 && time.Now().UnixNano() > item.expires {
		return "", errors.New("item with that key does not exist")
	}

	return item.data, nil
}

// GetMany ...
func (ml *MinimalismClient) GetMany(keys []string) (map[string]interface{}, []string, error) {
	var itemFound map[string]interface{}
	var itemNotFound []string

	for _, key := range keys {
		obj, err := ml.items.GetByKey(key)
		if err != nil {
			itemNotFound = append(itemNotFound, key)
		}

		item, ok := obj.(MinimalismItem)
		if !ok {
			return nil, nil, errors.New("can not map object to MinimalismItem model")
		}

		itemFound[key] = item.data
	}

	return itemFound, itemNotFound, nil
}

// Set ...
func (ml *MinimalismClient) Set(key string, value interface{}, expire time.Duration) error {
	var expires int64

	if expire > 0 {
		expires = time.Now().Add(expire).UnixNano()
	}

	if err := ml.items.Enqueue(key, MinimalismItem{
		data:    value,
		expires: expires,
	}); err != nil {
		return err
	}

	return nil
}

// Range ...
func (ml *MinimalismClient) Range(f func(key, value interface{}) bool) {
	now := time.Now().UnixNano()

	fn := func(key, value interface{}) bool {
		item := value.(MinimalismItem)

		if item.expires > 0 && now > item.expires {
			return true
		}

		return f(key, item.data)
	}

	ml.items.Range(fn)
}

// Delete deletes the key and its value from the cache.
func (ml *MinimalismClient) Delete(key string) error {
	return ml.items.Dequeue(key)
}

// GetNumberOfRecords return number of records
func (ml *MinimalismClient) GetNumberOfRecords() int {
	return ml.items.GetNumberOfKeys()
}

// GetDBSize method return redis database size
func (ml *MinimalismClient) GetCapacity() (interface{}, error) {
	return reflect.Type.Size(ml.items), nil
}

// Close closes the cache and frees up resources.
func (ml *MinimalismClient) Close() error {
	ml.close <- struct{}{}
	ml.items = queue.New(10 * 1024 * 1024) // 10 * 1024 * 1024 for 10 mb

	return nil
}
