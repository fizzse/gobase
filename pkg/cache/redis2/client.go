package redis2

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config 超时时间 单位秒
type Config struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DialTimeout  int    `yaml:"dialTimeout"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
	PoolSize     int    `yaml:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns"`
}

func NewClient(cfg *Config) (cli *redis.Client, cleanFunc func() error, err error) {
	ops := &redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     20,
		MinIdleConns: 10,
	}

	if cfg.DialTimeout != 0 {
		ops.DialTimeout = time.Duration(cfg.DialTimeout) * time.Second
	}

	if cfg.ReadTimeout != 0 {
		ops.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Second
	}

	if cfg.WriteTimeout != 0 {
		ops.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Second
	}

	if cfg.PoolSize != 0 {
		ops.PoolSize = cfg.PoolSize
	}

	if cfg.MinIdleConns != 0 {
		ops.MinIdleConns = cfg.MinIdleConns
	}

	rdb := redis.NewClient(ops)
	_, err = rdb.ClientID(context.Background()).Result() // 测试连通性
	return rdb, rdb.Close, err
}
