package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) IRedisCache {
	return &RedisCache{
		client: client,
	}
}

var ctx = context.Background()

func (c *RedisCache) Set(key string, value interface{}) error {
	tmp, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.client.Set(ctx, key, string(tmp), 10*time.Second).Result()
	return err
}

func (c *RedisCache) SetTTL(key string, value interface{}, ttl time.Duration) error {
	tmp, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.client.Set(ctx, key, string(tmp), ttl).Result()
	return err
}

func (c *RedisCache) Get(key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	} else {
		return value, nil
	}
}

func (c *RedisCache) Close() {
}

func (c *RedisCache) HSet(key string, field string, value interface{}) error {
	tmp, err := json.Marshal(value)
	if err != nil {
		return err
	}
	data := []interface{}{field, string(tmp)}
	_, err = c.client.HMSet(ctx, key, data).Result()
	return err
}

func (c *RedisCache) HDel(key string, field string) error {
	_, err := c.client.HDel(ctx, key, field).Result()
	return err
}

func (c *RedisCache) HGet(key string, field string) (string, error) {
	value, err := c.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", err
	}
	return value, err
}

func (c *RedisCache) HGetAll(key string) (map[string]string, error) {
	value, err := c.client.HGetAll(ctx, key).Result()
	return value, err
}

func (c *RedisCache) Del(key string) error {
	_, err := c.client.Del(ctx, key).Result()
	return err
}
