package caching

import (
	"time"

	"github.com/allegro/bigcache/v2"
)

// Config model for database caching config
type Config struct {
	CustomCache CustomCache     `json:"customCache,omitempty"`
	Redis       Redis           `json:"redis,omitempty"`
	BigCache    bigcache.Config `json:"bigCache,omitempty"`
}

// Redis model for redis config
type Redis struct {
	Password   string `json:"password,omitempty"`
	Host       string `json:"host,omitempty"`
	DB         int    `json:"db,omitempty"`
	MaxRetries int    `json:"maxRetries,omitempty"`
}

// CustomCache config model
type CustomCache struct {
	CacheSize        int64         `json:"cacheSize,omitempty"` // byte
	CleaningEnable   bool          `json:"cleaningEnable,omitempty"`
	CleaningInterval time.Duration `json:"cleaningInterval,omitempty"` // nanosecond
}

// customCacheItem private model for custom cache record
type customCacheItem struct {
	data    interface{} `json:"data,omitempty"`
	expires int64       `json:"expires,omitempty"`
}
