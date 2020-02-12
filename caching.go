package caching

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
)

// ICaching interface caching package
type ICaching interface {
	Middleware(hash hash.IHash) echo.MiddlewareFunc
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expire time.Duration) error
	Delete(key string) error
	GetCapacity() (interface{}, error)
	Close() error
}

const (
	MINIMALISM = iota
	BIGCACHE
	REDIS
)

// New function for Factory Pattern
func New(cachingType int, config *Config) ICaching {
	switch cachingType {
	case MINIMALISM:
		return NewMinimalism(config)
	case REDIS:
		return NewRedis(config)
	case BIGCACHE:
		return NewBigCache(config)
	}

	return nil
}
