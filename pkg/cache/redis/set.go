package redis

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

// set
func (c *Client) SAdd(key string, values ...interface{}) error {
	conn, err := c.GetConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	size := len(values) + 1
	fields := make([]interface{}, 0, size)
	fields = append(fields, key)
	for _, v := range values {
		fields = append(fields, v)
	}

	_, err = redis.Int64(conn.Do("SADD", fields...))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SRem(key string, value interface{}) error {
	conn, err := c.GetConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = redis.Int64(conn.Do("SREM", key, value))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SCard(key string) (int64, error) {
	conn, err := c.GetConn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	reply, err := redis.Int64(conn.Do("SCARD", key))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

func (c *Client) SMembers(key string) ([]string, error) {
	conn, err := c.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	reply, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}

	if len(reply) == 0 {
		return nil, errors.New("no data")
	}

	return reply, nil
}
