package redis

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

// add member
func (c *Client) ZAdd(key string, score int64, value interface{}) error {
	conn, err := c.GetConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.Int64(conn.Do("ZADD", key, score, value))
	if err != nil {
		return err
	}

	return nil
}

// range zset
func (c *Client) ZRange(key string, start, end int64, asc bool) ([]string, error) {
	conn, err := c.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 顺序
	if asc {
		reply, err := redis.Strings(conn.Do("ZRANGE", key, start, end))
		if err != nil {
			return nil, err
		}

		if len(reply) == 0 {
			return nil, errors.New("no data")
		}

		return reply, nil
	}

	// 逆序
	reply, err := redis.Strings(conn.Do("ZREVRANGE", key, start, end))
	if err != nil {
		return nil, err
	}

	if len(reply) == 0 {
		return nil, errors.New("no data")
	}

	return reply, nil
}

// range zset with score
func (c *Client) ZRangeWithScore(key string, start, end int64, asc bool) (map[string]int64, error) {
	conn, err := c.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if asc {
		reply, err := redis.Int64Map(conn.Do("ZRANGE", key, start, end, "WITHSCORES"))
		if err != nil {
			return nil, err
		}

		if len(reply) == 0 {
			return nil, errors.New("no data")
		}

		return reply, nil
	}

	reply, err := redis.Int64Map(conn.Do("ZREVRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	if len(reply) == 0 {
		return nil, errors.New("no data")
	}

	return reply, nil
}

func (c *Client) ZRangeByScoreWithScore(key string, min, max int64) (map[string]int64, error) {
	conn, err := c.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := redis.Int64Map(conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	if len(reply) == 0 {
		return nil, nil
	}

	return reply, nil
}

func (c *Client) ZRangeByScoreWithScoreAndLimit(key string, min, max, offset, limit int64) (map[string]int64, error) {
	conn, err := c.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := redis.Int64Map(conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES", "LIMIT", offset, limit))
	if err != nil {
		return nil, err
	}

	if len(reply) == 0 {
		return nil, nil
	}

	return reply, nil
}

func (c *Client) ZRem(key string, member string) (int64, error) {
	conn, err := c.GetConn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	reply, err := redis.Int64(conn.Do("ZREM", key, member))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

func (c *Client) ZRemRangeByScore(key string, min, max int64) (int64, error) {
	conn, err := c.GetConn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	reply, err := redis.Int64(conn.Do("ZREMRANGEBYSCORE", key, min, max))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

// query member score
func (c *Client) ZScore(key, member string) (int64, error) {
	conn, err := c.GetConn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	reply, err := redis.Int64(conn.Do("ZSCORE", key, member))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

// increment member score
func (c *Client) ZIncrBy(key, member string, increment int64) (int64, error) {
	conn, err := c.GetConn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	reply, err := redis.Int64(conn.Do("ZINCRBY", key, increment, member))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

// set len
func (c *Client) ZCard(key string) (int64, error) {
	conn, err := c.GetConn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	reply, err := redis.Int64(conn.Do("ZCARD", key))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

// member count in [min,max]
func (c *Client) ZCount(key string, min, max int64) (int64, error) {
	conn, err := c.GetConn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	reply, err := redis.Int64(conn.Do("ZCOUNT", key, min, max))
	if err != nil {
		return 0, err
	}

	return reply, nil
}
