package caching

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
)

// ICaching interface for this package
type ICaching interface {
	Middleware(hash hash.IHash) echo.MiddlewareFunc
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expire time.Duration) error
	Update(key string, value interface{}, expire time.Duration) error
	Delete(key string) error
	GetNumberOfRecords() int
	GetCapacity() (interface{}, error)
	Close() error
}

const (
	CUSTOM = iota
	BIGCACHE
	REDIS
)

// New instance based on caching type
func New(cachingType int, config *Config) ICaching {
	switch cachingType {
	case CUSTOM:
		return NewCustom(config)
	case REDIS:
		return NewRedis(config)
	case BIGCACHE:
		return NewBigCache(config)
	}

	return nil
}
