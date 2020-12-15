package redis

import "strconv"

/*
 * 计数器
 */

type Counter struct {
	entity     *Client // redis 实例
	expireTime int
}

func NewCounter(expireTime int, entity *Client) *Counter {
	if expireTime == 0 {
		expireTime = 7 * 24 * 3600 // 默认7天
	}

	client := new(Counter)
	client.expireTime = expireTime
	client.entity = entity

	return client
}

// 递增
func (c *Counter) INCR(key string) (int, error) {
	num, err := c.entity.INCR(key)
	if err != nil {
		return 0, err
	}

	// 设置过期时间
	if num == 1 {
		c.entity.Expire(key, c.expireTime)
	}

	return num, nil
}

func (c *Counter) GetScore(key string) (int, error) {
	ret, err := c.entity.GetString(key)
	if err != nil {
		return 0, err
	}

	num, err := strconv.Atoi(ret)
	if err != nil {
		return 0, err
	}

	return num, nil
}

// ZSet 相应 target 递增1
func (c *Counter) ZINCR(key, target string) (int, error) {
	num, err := c.entity.ZIncrBy(key, target, 1)
	if err != nil {
		return 0, err
	}

	// 设置过期时间
	if num == 1 {
		c.entity.Expire(key, c.expireTime)
	}
	return int(num), nil
}

// 查询 ZSet 中对应 target 的 score
func (c *Counter) GetZScore(key, target string) (int, error) {
	num, err := c.entity.ZScore(key, target)
	if err != nil {
		return 0, err
	}
	return int(num), nil
}
