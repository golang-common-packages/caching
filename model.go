package caching

import "github.com/allegro/bigcache/v2"

// Config model for database caching config
type Config struct {
	Redis    Redis           `json:"redis,omitempty"`
	BigCache bigcache.Config `json:"bigCache,omitempty"`
}

// Redis model provide info for redis config
type Redis struct {
	Prefix   string `json:"prefix,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	DB       int    `json:"db,omitempty"`
}
