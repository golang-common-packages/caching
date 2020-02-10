package caching

import (
	"log"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
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
func (bc *BigCacheClient) Set(key string, value string, expire time.Duration) (err error) {
	return bc.Client.Set(key, []byte(value))
}

// Get function will get value based on the key provided
func (bc *BigCacheClient) Get(key string) (value string, err error) {
	result, err := bc.Client.Get(key)
	value = string(result)
	return
}

// Delete function will delete value based on the key provided
func (bc *BigCacheClient) Delete(key string) (err error) {
	return bc.Client.Delete(key)
}

// Close function will close BigCache connection
func (bc *BigCacheClient) Close() {
	bc.Close()
}
