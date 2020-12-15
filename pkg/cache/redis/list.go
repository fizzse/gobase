package redis

/*
 * redis list 数据结构操作
 */

import (
	"github.com/gomodule/redigo/redis"
)

func (c *Client) LPush(key string, value interface{}) error {
	conn, err := c.GetConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = redis.Int64(conn.Do("LPUSH", key, value))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) LPop(key string) (string, error) {
	conn, err := c.GetConn()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	reply, err := redis.String(conn.Do("LPOP", key))
	if err != nil {
		return "", err
	}

	return reply, nil
}

func (c *Client) RPop(key string) (string, error) {
	conn, err := c.GetConn()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	reply, err := redis.String(conn.Do("RPOP", key))
	if err != nil {
		return "", err
	}

	return reply, nil
}

func (c *Client) BLPop(key string, timeout int) (string, error) {
	conn, err := c.GetConn()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	reply, err := redis.StringMap(conn.Do("BLPOP", key, timeout))
	if err != nil {
		return "", err
	}

	return reply[key], nil
}

func (c *Client) BRPop(key string, timeout int) (string, error) {
	conn, err := c.GetConn()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	reply, err := redis.StringMap(conn.Do("BRPOP", key, timeout))
	if err != nil {
		return "", err
	}

	return reply[key], nil
}

func (c *Client) LLength(key string) (int64, error) {
	conn, err := c.GetConn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	reply, err := redis.Int64(conn.Do("LLEN", key))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

func (c *Client) LRange(key string, start, end int) ([]string, error) {
	conn, err := c.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	reply, err := redis.Strings(conn.Do("LRANGE", key, start, end))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (c *Client) LRangeAll(key string) ([]string, error) {
	return c.LRange(key, 0, -1)
}

func (c *Client) LIndex(key string, index int) (string, error) {
	conn, err := c.GetConn()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	reply, err := redis.String(conn.Do("LINDEX", key, index))
	if err != nil {
		return "", err
	}

	return reply, nil
}

func (c *Client) LPushWithLimit(key string, value interface{}, limit int64) error {
	l, _ := c.LLength(key)
	if l >= limit {
		c.RPop(key)
	}

	return c.LPush(key, value)
}
