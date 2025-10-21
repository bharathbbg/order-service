package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bharathbbg/order-service/internal/config"
	"github.com/bharathbbg/order-service/internal/model"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(config config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func (c *RedisCache) CacheOrder(ctx context.Context, order *model.Order) error {
	key := fmt.Sprintf("order:%s", order.ID)
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, 24*time.Hour).Err()
}

func (c *RedisCache) GetCachedOrder(ctx context.Context, orderID string) (*model.Order, error) {
	key := fmt.Sprintf("order:%s", orderID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var order model.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
