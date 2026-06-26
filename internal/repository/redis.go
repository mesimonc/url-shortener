package repository

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

type Cache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewCache(redisURL string) (*Cache, error) {
    opts, err := redis.ParseURL(redisURL)
    if err != nil {
        return nil, err
    }

    client := redis.NewClient(opts)

    if err := client.Ping(context.Background()).Err(); err != nil {
        return nil, err
    }

    return &Cache{client: client, ttl: time.Hour}, nil
}

func (c *Cache) Get(ctx context.Context, code string) (string, error) {
    val, err := c.client.Get(ctx, "url:"+code).Result()
    if err == redis.Nil {
        return "", nil
    }
    if err != nil {
        return "", err
    }
    return val, nil
}

func (c *Cache) Set(ctx context.Context, code, originalURL string) error {
    return c.client.Set(ctx, "url:"+code, originalURL, c.ttl).Err()
}
