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

// Redis model provide info for redis config
type Redis struct {
	Password   string `json:"password,omitempty"`
	Host       string `json:"host,omitempty"`
	DB         int    `json:"db,omitempty"`
	MaxRetries int    `json:"maxRetries,omitempty"`
}

// CustomCache ...
type CustomCache struct {
	CleaningInterval time.Duration `json:"cleaningInterval,omitempty"` // nanosecond
	CacheSize        int64         `json:"cacheSize,omitempty"`        // byte
	SizeChecker      bool          `json:"sizeChecker,omitempty"`
}

// CustomCacheItem ...
type CustomCacheItem struct {
	data    interface{} `json:"data,omitempty"`
	expires int64       `json:"expires,omitempty"`
}
