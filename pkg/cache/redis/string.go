package redis

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

func (c *Client) SetStringEx(key, value string, expire int) error {
	conn, err := c.GetConn()
	if err != nil {
		return err
	}
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

func (c *Client) SetStringNx(key, value string, expire int) error {
	conn, err := c.GetConn()
	if err != nil {
		return err
	}
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

func (c *Client) SetStringNxPx(key, value string, expire int) error {
	conn, err := c.GetConn()
	if err != nil {
		return err
	}
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

func (c *Client) GetString(key string) (string, error) {
	conn, err := c.GetConn()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return reply, nil
}

func (c *Client) GetBytes(key string) ([]byte, error) {
	conn, err := c.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (c *Client) INCR(key string) (int, error) {
	conn, err := c.GetConn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	reply, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		return 0, err
	}

	return reply, nil
}
