package caching

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
)

// ICaching interface caching package
type ICaching interface {
	Middleware(hash hash.IHash) echo.MiddlewareFunc
	Get(key string) (string, error)
	Delete(key string) error
	Set(key string, value string, expire time.Duration) error
	GetCapacity() (result interface{}, err error)
	Close() error
}

const (
	REDIS = iota
	BIGCACHE
)

// New function for Factory Pattern
func New(cachingType int, config *Config) ICaching {
	switch cachingType {
	case REDIS:
		return NewRedis(config)
	case BIGCACHE:
		return NewBigCache(config)
	}

	return nil
}
