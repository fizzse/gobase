package delayqueue

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	qredis "github.com/fizzse/gobase/pkg/cache/redis"
)

/*
 * 基于redis实现,适用于响应时间没那么高的场景 1s
 */

type QueueByRedis struct {
	cfg          Config
	lastDealTime int64
	stopped      bool
	cli          *qredis.Client
}

func newQueueByRedis(config Config) (queue *QueueByRedis, err error) {
	redisCli, _, err := qredis.NewClient(&config.RedisCfg)
	if err != nil {
		return
	}

	queue = &QueueByRedis{cfg: config, cli: redisCli}
	return
}

func (q *QueueByRedis) Producer(msg string, expireTime int64) (err error) {
	pair := &redis.Z{Member: msg, Score: float64(expireTime)}
	err = q.cli.GetInterface().ZAdd(context.Background(), q.cfg.Name, pair).Err()
	return
}

func (q *QueueByRedis) Consumer(dealFunc DealFunc) {
	for {
		q.consumer(dealFunc)
		if q.stopped {
			break
		}

		time.Sleep(time.Second)
	}
}

func (q *QueueByRedis) consumer(dealFunc DealFunc) {
	t := time.Now().Unix()

	filter := &redis.ZRangeBy{Min: strconv.FormatInt(q.lastDealTime, 10), Max: strconv.FormatInt(t, 10)}
	msgs, err := q.cli.GetInterface().ZRangeByScore(context.Background(), q.cfg.Name, filter).Result()
	if err != nil {
		return
	}

	for _, msg := range msgs {
		if ierr := dealFunc(msg); ierr == nil {
			_ = q.cli.GetInterface().ZRem(context.Background(), q.cfg.Name, msg).Err()
		}
	}
	q.lastDealTime = t
	return
}

func (q *QueueByRedis) Stop() {
	q.stopped = true
}
