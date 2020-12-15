package redis

/*
 * redis 配置
 */

import (
	"fmt"
	"time"
)

type Config struct {
	Name        string        `json:"name" env:"REDIS_NAME"`
	Host        string        `json:"host" env:"REDIS_HOST"`
	Port        string        `json:"port" env:"REDIS_PORT"`
	Password    string        `json:"password" env:"REDIS_PWD"`
	MaxIdle     int           `json:"maxIdle" default:"50" env:"REDIS_MAX_IDLE"`
	MaxActive   int           `json:"maxActive" default:"100" env:"REDIS_MAX_ACTIVE"`
	IdleTimeout time.Duration `json:"idle_timeout"`
}

// RedisConfig builder pattern code
type ConfigBuilder struct {
	redisConfig *Config
}

func NewConfigBuilder() *ConfigBuilder {
	redisConfig := &Config{}
	b := &ConfigBuilder{redisConfig: redisConfig}
	return b
}

func (b *ConfigBuilder) Name(name string) *ConfigBuilder {
	b.redisConfig.Name = name
	return b
}

func (b *ConfigBuilder) Host(host string) *ConfigBuilder {
	b.redisConfig.Host = host
	return b
}

func (b *ConfigBuilder) Port(port string) *ConfigBuilder {
	b.redisConfig.Port = port
	return b
}

func (b *ConfigBuilder) Password(password string) *ConfigBuilder {
	b.redisConfig.Password = password
	return b
}

func (b *ConfigBuilder) MaxIdle(maxIdle int) *ConfigBuilder {
	b.redisConfig.MaxIdle = maxIdle
	return b
}

func (b *ConfigBuilder) MaxActive(maxActive int) *ConfigBuilder {
	b.redisConfig.MaxActive = maxActive
	return b
}

func (b *ConfigBuilder) IdleTimeout(idleTimeout time.Duration) *ConfigBuilder {
	b.redisConfig.IdleTimeout = idleTimeout
	return b
}

func (b *ConfigBuilder) Build() (*Config, error) {
	if b.redisConfig.Host == "" {
		return nil, fmt.Errorf("redis: host must be set")
	}

	if b.redisConfig.Port == "" {
		return nil, fmt.Errorf("redis: port must be set")
	}

	return b.redisConfig, nil
}
