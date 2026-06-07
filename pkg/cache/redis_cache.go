package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheService struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewRedisCacheService(client *redis.Client) *RedisCacheService {
	return &RedisCacheService{
		redisClient: client,
		ctx:         context.Background(),
	}
}

func (rs *RedisCacheService) Set(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)

	if err != nil {
		return err
	}
	return rs.redisClient.Set(rs.ctx, key, data, ttl*time.Minute).Err()
}

func (rs *RedisCacheService) Get(key string, dest any) error {
	val, err := rs.redisClient.Get(rs.ctx, key).Result()

	if err == redis.Nil {
		return redis.Nil
	}

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (rs *RedisCacheService) Clear(pattern string) error {
	cursor := uint64(0)

	for {
		key, nextCursor, err := rs.redisClient.Scan(rs.ctx, cursor, pattern, 5).Result()

		if err != nil {
			return err
		}
		if len(key) > 0 {
			_, err := rs.redisClient.Del(rs.ctx, key...).Result()
			if err != nil {
				return err
			}
			cursor = nextCursor
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}

func (rs *RedisCacheService) Push(ctx context.Context, queueKey string, value any) error {

	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return rs.redisClient.LPush(ctx, queueKey, data).Err()
}

func (rs *RedisCacheService) Pop(ctx context.Context, queueKey string, value any) error {
	data, err := rs.redisClient.BRPop(ctx, 0, queueKey).Result()

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data[1]), value)
}

// Publish gửi message lên một Redis Pub/Sub channel
func (rs *RedisCacheService) Publish(ctx context.Context, channel string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rs.redisClient.Publish(ctx, channel, data).Err()
}

// Subscribe lắng nghe message từ một Redis Pub/Sub channel
func (rs *RedisCacheService) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return rs.redisClient.Subscribe(ctx, channel)
}
