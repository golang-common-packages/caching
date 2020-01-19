package caching

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
)

// ICaching interface cachestore package
type ICaching interface {
	Middleware(hash hash.IHash) echo.MiddlewareFunc
	Get(key string) (string, error)
	Delete(key string) error
	Set(key string, value string, expire time.Duration) error
}

const (
	REDIS = iota
)

// New function for Factory Pattern
func New(datastoreType int, config *CachingConfig) ICaching {
	switch datastoreType {
	case REDIS:
		return NewRedis(config)
	}

	return nil
}
