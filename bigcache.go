package caching

import (
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

func (bc *BigCacheClient) Middleware(hash hash.IHash) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			key := hash.SHA512(token)

			if val, err := bc.Get(key); err != nil {
				log.Printf("Can not get accesstoken from redis in redis middleware: %s", err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			} else if val == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// Set function will set key and value
func (bc *BigCacheClient) Set(key string, value interface{}, expire time.Duration) error {
	b, ok := value.(*[]byte)
	if !ok {
		return errors.New("value must be []byte")
	}

	return bc.Client.Set(key, *b)
}

// GetByKey function will get value based on the key provided
func (bc *BigCacheClient) Get(key string) (interface{}, error) {
	return bc.Client.Get(key)
}

// Delete function will delete value based on the key provided
func (bc *BigCacheClient) Delete(key string) error {
	return bc.Client.Delete(key)
}

// GetDBSize method return redis database size
func (bc *BigCacheClient) GetCapacity() (interface{}, error) {
	return bc.Client.Capacity(), nil
}

// Close function will close BigCache connection
func (bc *BigCacheClient) Close() error {
	return bc.Client.Close()
}
