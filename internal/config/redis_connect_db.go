package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClientConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisClient() (*redis.Client, error) {

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		return nil, fmt.Errorf("connect redis failed with error: %v", err)
	}

	cfg := &RedisClientConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	_, errP := rdb.Ping(context.Background()).Result()

	if errP != nil {
		return nil, fmt.Errorf("connect redis failed: %w", err)
	}
	fmt.Println("[INFO] Connect resdis successfully")
	return rdb, nil
}
