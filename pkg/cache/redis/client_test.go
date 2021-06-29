package redis

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestNewSingleCli(t *testing.T) {
	cfg := &Config{
		Password:     "s",
		DialTimeout:  1,
		ReadTimeout:  1,
		WriteTimeout: 1,
		PoolSize:     20,
		MinIdleConns: 10,
	}

	cfg.Single.Addr = "127.0.0.1:6379"
	cfg.Mode = ModeSingle

	cli, cleanFunc, err := NewSingleCli(cfg)
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

func TestNewSentinelCli(t *testing.T) {
	cfg := &Config{
		Mode:         ModeSentinel,
		Password:     "s",
		DialTimeout:  1,
		ReadTimeout:  1,
		WriteTimeout: 1,
		PoolSize:     20,
		MinIdleConns: 10,
	}

	cfg.Sentinel.MasterName = "mymaster"
	cfg.Sentinel.Addrs = []string{"127.0.0.1:26379"}
	cli, cleanFunc, err := NewSentinelCli(cfg)
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

func TestNewClient(t *testing.T) {
	cfg := &Config{
		Password:     "s",
		DialTimeout:  1,
		ReadTimeout:  1,
		WriteTimeout: 1,
		PoolSize:     20,
		MinIdleConns: 10,
	}

	cfg.Single.Addr = "127.0.0.1:6379"
	cfg.Mode = ModeSingle

	cli, cleanFunc, err := NewClient(cfg)
	if err != nil {

	}

	defer cleanFunc()

	conn := cli.GetInterface()
	_, err = conn.SetEX(context.Background(), "hello", "world", 3*time.Second).Result()
	if err != nil {
		log.Fatal("set key failed", err)
	}

	res, err := conn.Get(context.Background(), "hello").Result()
	if err != nil {
		log.Fatal("get key failed", err)
	}

	log.Println("res", res)
}
