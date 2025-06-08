// Package redis оборачивает подключение к Redis.
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

// Config — параметры подключения
type Config struct {
	Addr     string
	Password string
	DB       int
}

// New создает подключение к Redis
func New(cfg Config) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		DialTimeout:  3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: connect ping error: %w", err)
	}

	return &Redis{Client: rdb}, nil
}

// Close закрывает соединение
func (r *Redis) Close() error {
	return r.Client.Close()
}
