// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/fizzse/gobase/internal/gobase/dao"
	"github.com/fizzse/gobase/internal/gobase/option"
	"github.com/fizzse/gobase/internal/gobase/server/consumer"
	"github.com/fizzse/gobase/internal/gobase/server/rest"
	"github.com/fizzse/gobase/internal/gobase/server/rpc"
	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
	"github.com/fizzse/gobase/pkg/logger"
	"github.com/fizzse/gobase/pkg/trace"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitApp() (*App, func(), error) {
	config := option.LoadRestConfig()
	dbConfig := option.LoadDbConfig()
	dbCtx, cleanup, err := db.NewConn(dbConfig)
	if err != nil {
		return nil, nil, err
	}
	redisConfig := option.LoadRedisConfig()
	client, cleanup2, err := redis.NewClient(redisConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	sampleDao, cleanup3, err := dao.NewInstance(dbCtx, client)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	loggerConfig := option.LoadLoggerConfig()
	sugaredLogger, err := logger.New(loggerConfig)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	sampleBiz, cleanup4, err := biz.NewInstance(sampleDao, sugaredLogger)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	server, err := rest.New(config, sampleBiz)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	rpcConfig := option.LoadGrpcConfig()
	rpcServer, cleanup5, err := rpc.New(rpcConfig, sampleBiz)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	kafkaConfig := option.LoadConsumerConfig()
	scheduler, err := consumer.NewScheduler(sugaredLogger, kafkaConfig, sampleBiz)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	traceConfig := option.LoadTraceConfig()
	tracer, cleanup6, err := trace.New(traceConfig)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	app, cleanup7, err := NewApp(server, rpcServer, scheduler, sugaredLogger, tracer)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	return app, func() {
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var (
	logProvider      = wire.NewSet(logger.New, option.LoadLoggerConfig)
	traceProvider    = wire.NewSet(trace.New, option.LoadTraceConfig)
	dbProvider       = wire.NewSet(db.NewConn, option.LoadDbConfig)
	redisProvider    = wire.NewSet(redis.NewClient, option.LoadRedisConfig)
	daoProvider      = wire.NewSet(dao.NewInstance)
	bizProvider      = wire.NewSet(biz.NewInstance)
	grpcProvider     = wire.NewSet(rpc.New, option.LoadGrpcConfig)
	restProvider     = wire.NewSet(rest.New, option.LoadRestConfig)
	consumerProvider = wire.NewSet(consumer.NewScheduler, option.LoadConsumerConfig)
)
