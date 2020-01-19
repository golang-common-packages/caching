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
	Prefix string
}

// NewRedis function return redis client based on singleton pattern
func NewRedis(config *CachingConfig) ICaching {
	currentSession := &RedisClient{nil, ""}
	client, err := currentSession.connect(config.Redis)
	if err != nil {
		panic(err)
	} else {
		currentSession.Client = client
		currentSession.Prefix = config.Redis.Prefix
		log.Println("Connected to Redis Server")
	}

	return currentSession
}

// connect private function establish redis connection
func (r *RedisClient) connect(data Redis) (client *redis.Client, err error) {
	if r.Client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     data.Host,
			Password: data.Password,
			DB:       data.DB,
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

// Middleware function will provide an echo middleware for Redis
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

// Set function will set key and value
func (r *RedisClient) Set(key string, value string, expire time.Duration) (err error) {
	err = r.Client.Set(r.Prefix+key, value, expire).Err()
	return
}

// Get function will get value based on the key provided
func (r *RedisClient) Get(key string) (value string, err error) {
	value, err = r.Client.Get(r.Prefix + key).Result()
	return
}

// Delete function will delete value based on the key provided
func (r *RedisClient) Delete(key string) (err error) {
	err = r.Client.Del(r.Prefix + key).Err()
	return
}

// Close function will close redis connection
func (r *RedisClient) Close() {
	r.Close()
}
