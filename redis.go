package caching

import (
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
)

// RedisClient manage all redis action
type RedisClient struct {
	Client *redis.Client
}

// NewRedis function return redis Client based on singleton pattern
func NewRedis(config *Config) ICaching {
	currentSession := &RedisClient{nil}
	client, err := currentSession.connect(config.Redis)
	if err != nil {
		panic(err)
	} else {
		currentSession.Client = client
		log.Println("Connected to Redis Server")
	}

	return currentSession
}

// connect private method establish redis connection
func (r *RedisClient) connect(data Redis) (client *redis.Client, err error) {
	if r.Client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:       data.Host,
			Password:   data.Password,
			DB:         data.DB,
			MaxRetries: data.MaxRetries,
		})

		_, err := client.Ping().Result()
		if err != nil {
			log.Println("Fail to connect redis: ", err)
			return nil, err
		}
	} else {
		client = r.Client
		err = nil
	}
	return
}

// Middleware method will provide an echo middleware for Redis
func (r *RedisClient) Middleware(hash hash.IHash) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			key := hash.SHA512(token)

			if val, err := r.Get(key); err != nil {
				log.Printf("Can not get accesstoken from redis in redis middleware: %s", err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			} else if val == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// Set method will set key and value
func (r *RedisClient) Set(key string, value interface{}, expire time.Duration) error {
	return r.Client.Set(key, value, expire).Err()
}

// GetByKey method return value based on the key provided
func (r *RedisClient) Get(key string) (interface{}, error) {
	return r.Client.Get(key).Result()
}

// Delete method delete value based on the key provided
func (r *RedisClient) Delete(key string) error {
	return r.Client.Del(key).Err()
}

// GetNumberOfRecords return number of records
func (r *RedisClient) GetNumberOfRecords() int {
	return len(r.Client.Do("KEYS", "*").Args())
}

// GetDBSize method return redis database size
func (r *RedisClient) GetCapacity() (interface{}, error) {
	IntCmd := r.Client.DBSize()

	return IntCmd.Result()
}

// Close method will close redis connection
func (r *RedisClient) Close() error {
	return r.Client.Close()
}
