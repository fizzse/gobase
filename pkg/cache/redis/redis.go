package redis

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

func newPool(config *RedisConfig) *redis.Pool {
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

func NewClient(config *RedisConfig) *RedisClient {
	client := &RedisClient{}
	client.pool = newPool(config)
	return client
}

type RedisClient struct {
	pool *redis.Pool
}

func (c *RedisClient) GetConn() redis.Conn {
	conn := c.pool.Get()
	if conn.Err() != nil {
		log.Println("RedisClient get conn failed: ", conn.Err())
	}

	return conn
}
