package redis2

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestNewSentinelClusterClient(t *testing.T) {
	cfg := &SentinelClusterConfig{
		MasterName:    "mymaster",
		SentinelAddrs: []string{"127.0.0.1:26379"},
		Password:      "s",
		DialTimeout:   1,
		ReadTimeout:   1,
		WriteTimeout:  1,
		PoolSize:      20,
		MinIdleConns:  10,
	}

	cli, cleanFunc, err := NewSentinelClusterClient(cfg)
	if err != nil {

	}

	defer cleanFunc()
	_, err = cli.SetEX(context.Background(), "hello", "world", 3*time.Second).Result()
	if err != nil {
		log.Fatal("set key failed", err)
	}

	res, err := cli.Get(context.Background(), "hello").Result()
	if err != nil {
		log.Fatal("get key failed", err)
	}

	log.Println("res", res)
}
