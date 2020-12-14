package redis

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

func (c *RedisClient) SetStringEx(key, value string, expire int) error {
	conn := c.GetConn()
	defer conn.Close()

	reply, err := redis.String(conn.Do("SET", key, value, "EX", expire))
	if err != nil {
		return err
	}
	if reply != "OK" {
		return errors.New("redis set key ex error")
	}

	return nil
}

func (c *RedisClient) SetStringNx(key, value string, expire int) error {
	conn := c.GetConn()
	defer conn.Close()

	reply, err := redis.String(conn.Do("SET", key, value, "NX", "EX", expire))
	if err != nil {
		return err
	}
	if reply != "OK" {
		return errors.New("redis set key nx error")
	}

	return nil
}

func (c *RedisClient) SetStringNxPx(key, value string, expire int) error {
	conn := c.GetConn()
	defer conn.Close()

	reply, err := redis.String(conn.Do("SET", key, value, "NX", "PX", expire))
	if err != nil {
		return err
	}
	if reply != "OK" {
		return errors.New("redis set key nx error")
	}

	return nil
}

func (c *RedisClient) GetString(key string) (string, error) {
	conn := c.GetConn()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return reply, nil
}

func (c *RedisClient) GetBytes(key string) ([]byte, error) {
	conn := c.GetConn()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (c *RedisClient) INCR(key string) (int, error) {
	conn := c.GetConn()
	defer conn.Close()

	reply, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		return 0, err
	}

	return reply, nil
}
