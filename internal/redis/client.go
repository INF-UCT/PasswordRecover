package redis

import (
	"context"
	"time"

	"chpassword/internal/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
	ttl time.Duration
}

func New(cfg config.RedisConfig, ttlMinutes int) (*Client, error) {
	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}

	if cfg.Password != "" {
		opts.Password = cfg.Password
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{
		rdb: rdb,
		ttl: time.Duration(ttlMinutes) * time.Minute,
	}, nil
}

func (c *Client) SaveToken(ctx context.Context, token, email string) error {
	return c.rdb.Set(ctx, token, email, c.ttl).Err()
}

func (c *Client) GetEmailByToken(ctx context.Context, token string) (string, error) {
	return c.rdb.Get(ctx, token).Result()
}

func (c *Client) DeleteToken(ctx context.Context, token string) error {
	return c.rdb.Del(ctx, token).Err()
}

func (c *Client) TokenExists(ctx context.Context, token string) (bool, error) {
	n, err := c.rdb.Exists(ctx, token).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (c *Client) GetTokenTTL(ctx context.Context, token string) (time.Duration, error) {
	return c.rdb.TTL(ctx, token).Result()
}
