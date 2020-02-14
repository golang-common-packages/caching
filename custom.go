package caching

import (
	"errors"
	"log"
	"net/http"
	"time"
	"unsafe"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
	"github.com/golang-common-packages/linear"
)

// CustomClient ...
type CustomClient struct {
	client *linear.Client
	close  chan struct{}
}

// NewCustom ...
func NewCustom(config *Config) ICaching {
	currentSession := &CustomClient{linear.New(config.CustomCache.CacheSize, config.CustomCache.SizeChecker), make(chan struct{})}

	// Check record expiration time and remove
	go func() {
		ticker := time.NewTicker(config.CustomCache.CleaningInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				items := currentSession.client.GetItems()
				items.Range(func(key, value interface{}) bool {
					item := value.(CustomCacheItem)

					if item.expires < time.Now().UnixNano() {
						k, _ := key.(string)
						currentSession.client.Get(k)
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
func (cl *CustomClient) Middleware(hash hash.IHash) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			key := hash.SHA512(token)

			if val, err := cl.Read(key); err != nil {
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
func (cl *CustomClient) Get(key string) (interface{}, error) {
	obj, err := cl.client.Get(key)
	if err != nil {
		return nil, err
	}

	item, ok := obj.(CustomCacheItem)
	if !ok {
		return nil, errors.New("can not map object to CustomCacheItem model")
	}

	if item.expires < time.Now().UnixNano() {
		return nil, nil
	}

	return item.data, nil
}

// GetMany ...
func (cl *CustomClient) GetMany(keys []string) (map[string]interface{}, []string, error) {
	var itemFound map[string]interface{}
	var itemNotFound []string

	for _, key := range keys {
		obj, err := cl.client.Get(key)
		if obj == nil && err == nil {
			itemNotFound = append(itemNotFound, key)
		}

		item, ok := obj.(CustomCacheItem)
		if !ok {
			return nil, nil, errors.New("can not map object to CustomCacheItem model")
		}

		itemFound[key] = item.data
	}

	return itemFound, itemNotFound, nil
}

// Read ...
func (cl *CustomClient) Read(key string) (interface{}, error) {
	obj, err := cl.client.Read(key)
	if err != nil {
		return nil, err
	}

	item, ok := obj.(CustomCacheItem)
	if !ok {
		return nil, errors.New("can not map object to CustomCacheItem model in Read method")
	}

	if item.expires < time.Now().UnixNano() {
		return nil, nil
	}

	return item.data, nil
}

// Set ...
func (cl *CustomClient) Set(key string, value interface{}, expire time.Duration) error {
	if err := cl.client.Push(key, CustomCacheItem{
		data:    value,
		expires: time.Now().Add(expire).UnixNano(),
	}); err != nil {
		return err
	}

	return nil
}

func (cl *CustomClient) Update(key string, value interface{}) error {
	return cl.client.Update(key, value)
}

// Delete deletes the key and its value from the cache.
func (cl *CustomClient) Delete(key string) error {
	_, err := cl.client.Get(key)
	return err
}

// Range ...
func (cl *CustomClient) Range(f func(key, value interface{}) bool) {
	now := time.Now().UnixNano()

	fn := func(key, value interface{}) bool {
		item := value.(CustomCacheItem)

		if item.expires > 0 && now > item.expires {
			return true
		}

		return f(key, item.data)
	}

	cl.client.Range(fn)
}

// GetNumberOfRecords return number of records
func (cl *CustomClient) GetNumberOfRecords() int {
	return cl.client.GetNumberOfKeys()
}

// GetDBSize method return redis database size
func (cl *CustomClient) GetCapacity() (interface{}, error) {
	return unsafe.Sizeof(cl.client), nil
}

// Close closes the cache and frees up resources.
func (cl *CustomClient) Close() error {
	cl.close <- struct{}{}
	cl.client = linear.New(10*1024*1024, true) // 10 * 1024 * 1024 for 10 mb

	return nil
}