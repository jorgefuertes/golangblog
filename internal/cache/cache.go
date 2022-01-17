package cache

import (
	"sync"
	"time"

	"golangblog/internal/log"

	"github.com/go-redis/redis"
)

// 1 month
const OneMonth = time.Duration(24 * 60 * 60 * 30 * time.Second)
const OneYear = time.Duration(24 * 60 * 60 * 365 * time.Second)

// thread safe instance
var connOnce sync.Once
var c *redis.Client

// Begin the connection
func Begin() {
	connOnce.Do(func() {
		var err *redis.StatusCmd
		for {
			c = redis.NewClient(&redis.Options{
				Addr: "localhost",
				DB:   0,
			})

			// check
			err = c.Set("client-test", time.Now(), 3*time.Second)
			if err.Err() == nil {
				log.Info("REDIS:CONN", "Connected")
				break
			}

			log.Error("REDIS:ERROR", err.Err())
			log.Info("REDIS:RETRY", "Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	})
}

// Set
func Set(key string, value interface{}) error {
	log.Debug("REDIS:SET", key, value)
	return c.Set(key, value, 0).Err()
}

// SetWithEx
func SetWithEx(key string, value interface{}, expire time.Duration) error {
	log.Debug("REDIS:SETEX", key, value, expire)
	return c.Set(key, value, expire).Err()
}

// Get
func Get(key string) string {
	log.Debug("REDIS:GET", key, c.Get(key).Val())
	return c.Get(key).Val()
}

func GetInt(key string) int {
	i, _ := c.Get(key).Int()
	return i
}

func GetBytes(key string) []byte {
	b, _ := c.Get(key).Bytes()
	return b
}

// Delete
func Del(key string) {
	c.Del(key)
}

// Exists
func Exists(key string) bool {
	return (c.Exists(key).Val() == 1)
}
