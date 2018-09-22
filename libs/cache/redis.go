package cache

import (
	"sync"

	"github.com/garyburd/redigo/redis"
	"github.com/kere/gno/libs/conf"
	libRedis "github.com/kere/gno/libs/redis"
)

// RedisCache redis
type RedisCache struct {
	Driver string
	client *libRedis.Pool
	mutex  sync.RWMutex
}

// NewRedisCache new
func NewRedisCache() *RedisCache {
	return &RedisCache{Driver: "redis"}
}

// GetDriver string
func (r *RedisCache) GetDriver() string {
	return r.Driver
}

// Init err
func (r *RedisCache) Init(c conf.Conf) error {
	r.client = libRedis.NewPool(c)
	return nil
}

// GetRedis return redis client
func (r *RedisCache) GetRedis() *libRedis.Pool {
	return r.client
}

// Delete remove key
func (r *RedisCache) Delete(key string) error {
	return r.client.Delete(key)
}

// Set func
func (r *RedisCache) Set(key, value string, expire int) error {
	var err error

	r.mutex.Lock()
	if expire > 0 {
		err = r.client.Send("SETEX", key, expire, value)
	} else {
		err = r.client.Send("SET", key, value)
	}
	r.mutex.Unlock()
	if err != nil {
		return err
	}
	return nil
}

// Get func
func (r *RedisCache) Get(key string) (interface{}, error) {
	// r.mutex.RLock()
	return r.client.DoString("GET", key)
	// r.mutex.RUnlock()
}

// Exists key
func (r *RedisCache) IsExists(key string) (bool, error) {
	return r.client.Exists(key)
}

// GetString string
func (r *RedisCache) GetString(key string) (string, error) {
	return redis.String(r.Get(key))
}

// GetInt int
func (r *RedisCache) GetInt(key string) (int, error) {
	return redis.Int(r.Get(key))
}

// GetInt64 int64
func (r *RedisCache) GetInt64(key string) (int64, error) {
	return redis.Int64(r.Get(key))
}

// GetFloat float64
func (r *RedisCache) GetFloat(key string) (float64, error) {
	return redis.Float64(r.Get(key))
}
