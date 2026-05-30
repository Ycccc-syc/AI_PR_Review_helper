package repository

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis(addr, password string) (*redis.Client, error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}
	return rdb, nil
}

func GetRedis() *redis.Client {
	return rdb
}

// CacheKey 基于 PR 标识 + diff 内容哈希生成缓存键，新 commit 导致新键、旧缓存自然淘汰
func CacheKey(owner, repo string, number int, diff string) string {
	h := sha256.Sum256([]byte(diff))
	return fmt.Sprintf("review:cache:%s:%s:%d:%x", owner, repo, number, h[:8])
}

func GetCachedResult(key string) ([]byte, bool) {
	if rdb == nil {
		return nil, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	data, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return data, true
}

func SetCachedResult(key string, value interface{}, ttl time.Duration) {
	if rdb == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	rdb.Set(ctx, key, data, ttl)
}
