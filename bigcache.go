package caching

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/allegro/bigcache/v2"
	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
)

// BigCacheClient manage all BigCache action
type BigCacheClient struct {
	Client *bigcache.BigCache
}

// NewBigCache init new instance
func NewBigCache(config *Config) ICaching {
	currentSession := &BigCacheClient{nil}
	client, err := bigcache.NewBigCache(config.BigCache)
	if err != nil {
		panic(err)
	} else {
		currentSession.Client = client
		log.Println("BigCache is ready")
	}

	return currentSession
}

// Middleware for echo framework
func (bc *BigCacheClient) Middleware(hash hash.IHash) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			key := hash.SHA512(token)

			if _, err := bc.Get(key); (err != nil) && (err.Error() != "Entry not found") {
				log.Printf("Can not get accesstoken from redis in redis middleware: %s", err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			} else {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// Set function will set key and value
func (bc *BigCacheClient) Set(key string, value interface{}, expire time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return errors.New("can not marshal value to []byte")
	}

	return bc.Client.Set(key, b)
}

// GetByKey function will get value based on the key provided
func (bc *BigCacheClient) Get(key string) (interface{}, error) {
	b, err := bc.Client.Get(key)
	if err != nil {
		return nil, err
	}

	var value interface{}
	json.Unmarshal(b, value)

	return value, nil
}

// Delete function will delete value based on the key provided
func (bc *BigCacheClient) Delete(key string) error {
	return bc.Client.Delete(key)
}

// GetNumberOfRecords return number of records
func (bc *BigCacheClient) GetNumberOfRecords() int {
	return bc.Client.Len()
}

// GetDBSize method return redis database size
func (bc *BigCacheClient) GetCapacity() (interface{}, error) {
	return bc.Client.Capacity(), nil
}

// Close function will close BigCache connection
func (bc *BigCacheClient) Close() error {
	return bc.Client.Close()
}
