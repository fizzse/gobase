package loader

import (
	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
)

/*
 * 加载配置
 */

type RestConfig struct {
	Host       string `json:"host" yaml:"host"`
	Port       int    `json:"port" yaml:"port"`
	DebugModel bool   `json:"debugModel" yaml:"debugModel"`
}

func LoadRestConfig() *RestConfig {
	return &RestConfig{
		Host: "0.0.0.0",
		Port: 8080,
	}
}

func LoadDbConfig() *db.Config {
	return &db.Config{
		Drive:    "mysql",
		Address:  "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "s",
		DbName:   "test",
		Charset:  "utf8",
	}
}

func LoadRedisConfig() *redis.Config {
	return &redis.Config{
		Host:     "127.0.0.1",
		Port:     "6379",
		Password: "s",
	}
}
