package redis

/*
 * redis 通用命令
 */

import "github.com/gomodule/redigo/redis"

func (c *Client) Del(key string) error {
	conn, err := c.GetConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.Int64(conn.Do("DEL", key))
	if err != nil {
		return err
	}

	return nil
}

// key 加过期时间
func (c *Client) Expire(key string, duration int) error {
	conn, err := c.GetConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.Int(conn.Do("EXPIRE", key, duration))
	if err != nil {
		return err
	}

	return nil
}

// 重命名key
func (c *Client) Rename(oldKey, newKey string) error {
	conn, err := c.GetConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.String(conn.Do("RENAME", oldKey, newKey))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Keys(keys string) ([]string, error) {
	conn, err := c.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res, err := redis.Strings(conn.Do("KEYS", keys))
	if err != nil {
		return nil, err
	}

	return res, nil
}
