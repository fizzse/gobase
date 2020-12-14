package redis

/*
 * redis hash 数据结构操作
 */

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

func (c *RedisClient) HGetAll(key string) (map[string]string, error) {
	conn := c.GetConn()
	defer conn.Close()

	v, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, errors.New("no data")
	}

	return v, nil
}

func (c *RedisClient) HGet(key string, field string) (string, error) {
	conn := c.GetConn()
	defer conn.Close()

	v, err := redis.String(conn.Do("HGET", key, field))
	if err != nil {
		return "", err
	}

	return v, nil
}

func (c *RedisClient) HSet(key string, fieldKey string, fieldValue interface{}) error {
	conn := c.GetConn()
	defer conn.Close()

	_, err := redis.Int(conn.Do("HSET", fieldKey, fieldValue))
	return err
}

func (c *RedisClient) HMSet(key string, fields map[string]interface{}) error {
	conn := c.GetConn()
	defer conn.Close()

	size := len(fields)
	args := make([]interface{}, 0, 2*size+1)
	args = append(args, key)
	for k, v := range fields {
		args = append(args, k)
		args = append(args, v)
	}

	_, err := redis.String(conn.Do("HMSET", args...))
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) HDel(key string, field string) error {
	conn := c.GetConn()
	defer conn.Close()

	_, err := redis.Int(conn.Do("HDEL", key, field))
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) HKeys(key string) ([]string, error) {
	conn := c.GetConn()
	defer conn.Close()

	fieldKeys, err := redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		return nil, err
	}

	return fieldKeys, nil
}
