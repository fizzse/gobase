package main

import (
	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
)

func main(){
	redisConfig:=&redis.Config{}
	redis.NewClient(redisConfig)

	dbConfig:= &db.Config{}
	db.NewConn(dbConfig)
}
