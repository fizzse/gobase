package delayqueue

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/fizzse/gobase/pkg/cache/redis"
)

func newRedisQueue() *QueueByRedis {
	redisCfg := redis.Config{Mode: redis.ModeSingle}
	redisCfg.Single.Addr = "127.0.0.1:6379"
	redisCfg.Password = "s"
	cfg := Config{
		Name:     "delayQueue:fizzse",
		RedisCfg: redisCfg,
	}

	queue, err := newQueueByRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return queue
}

func TestQueueByRedis_Producer(t *testing.T) {
	queue := newRedisQueue()
	timeNow := time.Now().Unix()
	queue.Producer("00001", timeNow+5)
	queue.Producer("00002", timeNow+6)
	queue.Producer("00003", timeNow+7)
	queue.Producer("00004", timeNow+7)
}

func TestQueueByRedis_Consumer(t *testing.T) {
	queue := newRedisQueue()
	go func() {
		queue.Consumer(func(msg string) error {
			fmt.Println(msg, time.Now())
			return nil
		})
	}()

	time.Sleep(10 * time.Second)
	queue.Stop()
}
