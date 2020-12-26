package server

import (
	"github.com/fizzse/gobase/internal/gobase/server/rest"
	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
)

func LoadRestConfig() *rest.Config {
	return &rest.Config{
		Host:       "0.0.0.0",
		Port:       8080,
		DebugModel: true,
	}
}

func LoadDbConfig() *db.Config {
	return &db.Config{
		Drive:    "mysql",
		Address:  "172.28.47.6",
		Port:     3306,
		User:     "root",
		Password: "s",
		DbName:   "gobase",
		Charset:  "utf8",
	}
}

func LoadRedisConfig() *redis.Config {
	return &redis.Config{
		Host:     "172.28.47.6",
		Port:     "6379",
		Password: "s",
	}
}
