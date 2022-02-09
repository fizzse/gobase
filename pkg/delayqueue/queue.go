package delayqueue

import "github.com/fizzse/gobase/pkg/cache/redis"

type DealFunc func(msg string) error

type DelayQueue interface {
	Producer(msg string, expireTime int64) error
	Consumer(dealFunc DealFunc)
}

type Config struct {
	Drive    string       `yaml:"drive"`    // 驱动 redis
	Name     string       `yaml:"name"`     // 队列名字
	RedisCfg redis.Config `yaml:"redisCfg"` //
}

func New(config Config) (queue DelayQueue, err error) {
	return newQueueByRedis(config)
}
