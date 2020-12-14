package redis

/*
 * redis 通用命令
 */

import "github.com/gomodule/redigo/redis"

func (c *RedisClient) Del(key string) error {
	conn := c.GetConn()
	defer conn.Close()
	_, err := redis.Int64(conn.Do("DEL", key))
	if err != nil {
		return err
	}

	return nil
}

// key 加过期时间
func (c *RedisClient) Expire(key string, duration int) error {
	conn := c.GetConn()
	defer conn.Close()

	_, err := redis.Int(conn.Do("EXPIRE", key, duration))
	if err != nil {
		return err
	}

	return nil
}

// 重命名key
func (c *RedisClient) Rename(oldKey, newKey string) error {
	conn := c.GetConn()
	defer conn.Close()

	_, err := redis.String(conn.Do("RENAME", oldKey, newKey))
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) Keys(keys string) ([]string, error) {
	conn := c.GetConn()
	defer conn.Close()

	res, err := redis.Strings(conn.Do("KEYS", keys))
	if err != nil {
		return nil, err
	}

	return res, nil
}
