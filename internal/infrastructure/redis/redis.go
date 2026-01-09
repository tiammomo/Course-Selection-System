package redis

import (
	"context"
	"fmt"
	"time"

	"course_select/internal/config"

	"github.com/gomodule/redigo/redis"
)

// Client Redis 客户端
type Client struct {
	pool *redis.Pool
}

// New 创建 Redis 客户端
func New(cfg *config.RedisConfig) (*Client, error) {
	pool := &redis.Pool{
		MaxActive:   cfg.PoolSize,
		MaxIdle:     cfg.PoolSize,
		IdleTimeout: time.Minute * 5,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.Addr(),
				redis.DialPassword(cfg.Password),
				redis.DialDatabase(cfg.DB),
				redis.DialConnectTimeout(cfg.DialTimeout),
				redis.DialReadTimeout(cfg.ReadTimeout),
				redis.DialWriteTimeout(cfg.WriteTimeout),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to connect redis: %w", err)
			}
			return c, nil
		},
	}

	// 测试连接
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &Client{pool: pool}, nil
}

// Close 关闭连接池
func (c *Client) Close() error {
	return c.pool.Close()
}

// Pool 获取连接池
func (c *Client) Pool() *redis.Pool {
	return c.pool
}

// Hash 操作
func (c *Client) HIncrBy(ctx context.Context, key string, field string, delta int) (int, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	result, err := conn.Do("HINCRBY", key, field, delta)
	if err != nil {
		return 0, err
	}
	return redis.Int(result, err)
}

// SIsMember 判断是否存在
func (c *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	result, err := conn.Do("SISMEMBER", key, member)
	if err != nil {
		return false, err
	}
	return redis.Bool(result, err)
}

// SMembers 获取所有成员
func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	result, err := conn.Do("SMEMBERS", key)
	if err != nil {
		return nil, err
	}
	return redis.Strings(result, err)
}

// LPush 左推入队列
func (c *Client) LPush(ctx context.Context, key string, values ...interface{}) (int, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	args := append([]interface{}{key}, values...)
	result, err := conn.Do("LPUSH", args...)
	if err != nil {
		return 0, err
	}
	return redis.Int(result, err)
}
