package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func newPool(config *Config) *redis.Pool {
	server := config.Host + ":" + config.Port
	return &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: config.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if config.Password != "" {
				if _, err := c.Do("AUTH", config.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func NewClient(config *Config) (redisCtx *Client, clean func(), err error) {
	redisCli := &Client{}
	redisCli.pool = newPool(config)

	// 初始化获取一次链接 来检查redis链接是否正常
	_, err = redisCli.GetConn()
	return redisCli, redisCli.Close, err
}

type Client struct {
	pool *redis.Pool
}

func (c *Client) GetConn() (redis.Conn, error) {
	conn := c.pool.Get()
	return conn, conn.Err()
}

func (c *Client) Close() {
	c.pool.Close()
}
