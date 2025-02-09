package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *RedisCache
	once        sync.Once
)

const (
	maxRetries = 3
	baseDelay  = 2 * time.Second
)

// RedisCache wraps the Redis client with additional methods.
type RedisCache struct {
	client *redis.Client
}

// GetInstance returns the singleton instance of RedisCache.
func GetInstance() *RedisCache {
	once.Do(func() {
		client, err := connectWithRetry()
		if err != nil {
			panic(fmt.Sprintf("Failed to initialize Redis client: %v", err))
		}
		redisClient = &RedisCache{client: client}
	})
	return redisClient
}

func connectWithRetry() (*redis.Client, error) {
	options := &redis.Options{
		Addr:     "localhost:6379", // Update with your Redis configuration.
		Password: "",
		DB:       0,
	}

	var client *redis.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		client = redis.NewClient(options)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err = client.Ping(ctx).Result(); err == nil {
			return client, nil
		}

		time.Sleep(time.Duration(i+1) * baseDelay)
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
}

// Exists checks if a key exists in Redis.
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error checking existence: %w", err)
	}
	return result > 0, nil
}

// Get retrieves and deserializes a value from Redis.
func (rc *RedisCache) Get(ctx context.Context, key string, target interface{}) error {
	data, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("key %s not found", key)
		}
		return fmt.Errorf("error getting value: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("error unmarshaling data: %w", err)
	}

	return nil
}

// Save serializes and stores a value in Redis.
func (rc *RedisCache) Save(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	if err := rc.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("error saving value: %w", err)
	}

	return nil
}

// Push appends a value to a Redis list.
func (rc *RedisCache) Push(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	if err := rc.client.RPush(ctx, key, data).Err(); err != nil {
		return fmt.Errorf("error pushing to list: %w", err)
	}

	return nil
}

// XAdd adds an entry to a Redis stream. This method guarantees that the message
// will be persisted until it is consumed.
func (rc *RedisCache) XAdd(ctx context.Context, args *redis.XAddArgs) (string, error) {
	return rc.client.XAdd(ctx, args).Result()
}

// RunScript runs the given Lua script with the specified keys and arguments.
func (rc *RedisCache) RunScript(ctx context.Context, script *redis.Script, keys []string, args ...interface{}) (interface{}, error) {
	return script.Run(ctx, rc.client, keys, args...).Result()
}

// IncrBy increments a key by the specified value.
func (rc *RedisCache) IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	return rc.client.IncrBy(ctx, key, increment).Result()
}

// Close closes the Redis connection.
func (rc *RedisCache) Close() error {
	if rc.client != nil {
		return rc.client.Close()
	}
	return nil
}
