// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package server

import (
	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/fizzse/gobase/internal/gobase/dao"
	"github.com/fizzse/gobase/internal/gobase/server/rest"
	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitApp() (*App, func(), error) {
	config := LoadRestConfig()
	dbConfig := LoadDbConfig()
	dbCtx, cleanup, err := db.NewConn(dbConfig)
	if err != nil {
		return nil, nil, err
	}
	redisConfig := LoadRedisConfig()
	client, cleanup2, err := redis.NewClient(redisConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	daoDao, cleanup3, err := dao.New(dbCtx, client)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	bizBiz, cleanup4, err := biz.New(daoDao)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	server, err := rest.New(config, bizBiz)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	app, cleanup5, err := NewApp(server)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	return app, func() {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var (
	daoProvider  = wire.NewSet(dao.New, db.NewConn, redis.NewClient, LoadDbConfig, LoadRedisConfig)
	bizProvider  = wire.NewSet(biz.New)
	restProvider = wire.NewSet(rest.New, LoadRestConfig)
)
