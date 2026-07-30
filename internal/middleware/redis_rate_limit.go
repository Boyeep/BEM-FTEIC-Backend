package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimitStore struct {
	client *redis.Client
	prefix string
}

func NewRedisRateLimitStore(rawURL, prefix string) (*RedisRateLimitStore, error) {
	options, err := redis.ParseURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse REDIS_URL: %w", err)
	}
	client := redis.NewClient(options)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}
	return &RedisRateLimitStore{client: client, prefix: prefix}, nil
}

func (s *RedisRateLimitStore) Increment(key string, now time.Time, window time.Duration) (int, time.Time, error) {
	ctx := context.Background()
	redisKey := s.prefix + ":" + key
	result, err := s.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, redisKey)
		pipe.ExpireNX(ctx, redisKey, window)
		return nil
	})
	if err != nil {
		return 0, time.Time{}, err
	}
	count, err := result[0].(*redis.IntCmd).Result()
	if err != nil {
		return 0, time.Time{}, err
	}
	ttl, err := s.client.TTL(ctx, redisKey).Result()
	if err != nil {
		return 0, time.Time{}, err
	}
	if ttl <= 0 {
		ttl = window
	}
	return int(count), now.Add(ttl), nil
}
