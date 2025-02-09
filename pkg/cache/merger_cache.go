package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// MergerCache is a wrapper that contains the RedisCache instance and the notification channel (implemented as a Redis stream).
type MergerCache struct {
	Cache   *RedisCache
	Channel string
}

// NewMergerCache is a factory method to create a new MergerCache with the given channel.
func NewMergerCache(channel string) *MergerCache {
	return &MergerCache{
		Cache:   GetInstance(),
		Channel: channel,
	}
}

// decrAndNotifyScript is a Lua script that atomically decrements the counter stored at a given key.
// If the counter reaches zero, it adds an entry to the Redis stream (the notification channel)
// containing the notification message and a timestamp.
var decrAndNotifyScript = redis.NewScript(`
    local count = redis.call('DECR', KEYS[1])
    if count == 0 then
        redis.call('XADD', ARGV[1], '*', 'message', ARGV[2], 'timestamp', ARGV[3])
    end
    return count
`)

// DecrementCounter atomically decrements the counter stored at counterKey.
// When the counter reaches 0, it sends a notification (using the provided notificationMessage)
// to the stream defined in the MergerCache's Channel field.
func (mc *MergerCache) DecrementCounter(ctx context.Context, counterKey string, notificationMessage string) (int64, error) {
	now := time.Now().Unix()
	result, err := mc.Cache.RunScript(ctx, decrAndNotifyScript, []string{counterKey}, mc.Channel, notificationMessage, now)
	if err != nil {
		return 0, fmt.Errorf("failed to run Lua script: %w", err)
	}

	count, ok := result.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected result type from Lua script")
	}

	return count, nil
}

// IncrementCounter increments the counter stored at counterKey by the specified amount.
func (mc *MergerCache) IncrementCounter(ctx context.Context, counterKey string, increment int64) (int64, error) {
	return mc.Cache.IncrBy(ctx, counterKey, increment)
}
