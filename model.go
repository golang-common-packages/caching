package caching

// CachingConfig model for database caching config
type CachingConfig struct {
	Redis Redis `json:"redis,omitempty"`
}

// Redis model provide info for redis config
type Redis struct {
	Prefix   string `json:"prefix,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	DB       int    `json:"db,omitempty"`
}
