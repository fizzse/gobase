package redis2

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	ModeSingle   = "single"
	ModeSentinel = "sentinel"
)

// Config 超时时间 单位秒
type Config struct {
	Mode   string `yaml:"mode"`
	Single struct {
		Addr string `yaml:"addr"`
	} `json:"single"`

	Sentinel struct {
		MasterName string   `yaml:"masterName"`
		Addrs      []string `yaml:"Addrs"`
	} `yaml:"sentinel"`

	Password     string        `yaml:"password"`
	DialTimeout  time.Duration `yaml:"dialTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	PoolSize     int           `yaml:"poolSize"`
	MinIdleConns int           `yaml:"minIdleConns"`
}

func NewClient(cfg *Config) (cli *Client, cleanFunc func(), err error) {
	cli = &Client{}
	switch cfg.Mode {
	case ModeSingle:
		cli.singleCli, cleanFunc, err = NewSingleCli(cfg)
	case ModeSentinel:
		cli.sentinelCli, cleanFunc, err = NewSentinelCli(cfg)
	}

	return
}

type Client struct {
	singleCli   *redis.Client
	sentinelCli *redis.ClusterClient
}

func (c *Client) GetSingleCli() *redis.Client {
	return c.singleCli
}

func (c *Client) GetSentinelCli() *redis.ClusterClient {
	return c.sentinelCli
}

// GetInterface 返回接口
func (c *Client) GetInterface() redis.Cmdable {
	if c.singleCli != nil {
		return c.singleCli
	}

	if c.sentinelCli != nil {
		return c.sentinelCli
	}

	return nil
}

func NewSingleCli(cfg *Config) (cli *redis.Client, cleanFunc func(), err error) {
	ops := &redis.Options{
		Addr:         cfg.Single.Addr,
		Password:     cfg.Password,
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     20,
		MinIdleConns: 10,
	}

	if cfg.DialTimeout != 0 {
		ops.DialTimeout = cfg.DialTimeout
	}

	if cfg.ReadTimeout != 0 {
		ops.ReadTimeout = cfg.ReadTimeout
	}

	if cfg.WriteTimeout != 0 {
		ops.WriteTimeout = cfg.WriteTimeout
	}

	if cfg.PoolSize != 0 {
		ops.PoolSize = cfg.PoolSize
	}

	if cfg.MinIdleConns != 0 {
		ops.MinIdleConns = cfg.MinIdleConns
	}

	rdb := redis.NewClient(ops)
	_, err = rdb.ClientID(context.Background()).Result() // 测试连通性
	return rdb, func() {
		rdb.Close()
	}, err
}

func NewSentinelCli(cfg *Config) (cli *redis.ClusterClient, cleanFunc func(), err error) {
	ops := &redis.FailoverOptions{
		MasterName:    cfg.Sentinel.MasterName,
		SentinelAddrs: cfg.Sentinel.Addrs,
		Password:      cfg.Password,
		DialTimeout:   time.Second,
		ReadTimeout:   time.Second,
		WriteTimeout:  time.Second,
		PoolSize:      20,
		MinIdleConns:  10,
	}

	if cfg.DialTimeout != 0 {
		ops.DialTimeout = cfg.DialTimeout
	}

	if cfg.ReadTimeout != 0 {
		ops.ReadTimeout = cfg.ReadTimeout
	}

	if cfg.WriteTimeout != 0 {
		ops.WriteTimeout = cfg.WriteTimeout
	}

	if cfg.PoolSize != 0 {
		ops.PoolSize = cfg.PoolSize
	}

	if cfg.MinIdleConns != 0 {
		ops.MinIdleConns = cfg.MinIdleConns
	}

	rdb := redis.NewFailoverClusterClient(ops)
	_, err = rdb.ClientID(context.Background()).Result() // 测试连通性
	return rdb, func() {
		rdb.Close()
	}, err
}
