// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package server

import (
	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/fizzse/gobase/internal/gobase/dao"
	"github.com/fizzse/gobase/internal/gobase/server/consumer"
	"github.com/fizzse/gobase/internal/gobase/server/rest"
	"github.com/fizzse/gobase/internal/gobase/server/rpc"
	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
	"github.com/fizzse/gobase/pkg/logger"
	"github.com/fizzse/gobase/pkg/trace"
	"github.com/google/wire"
)

var (
	logProvider      = wire.NewSet(logger.New, LoadLoggerConfig)
	traceProvider    = wire.NewSet(trace.New, LoadTraceConfig)
	dbProvider       = wire.NewSet(db.NewConn, LoadDbConfig)
	redisProvider    = wire.NewSet(redis.NewClient, LoadRedisConfig)
	daoProvider      = wire.NewSet(dao.NewInstance)
	bizProvider      = wire.NewSet(biz.NewInstance)
	grpcProvider     = wire.NewSet(rpc.New, LoadGrpcConfig)
	restProvider     = wire.NewSet(rest.New, LoadRestConfig)
	consumerProvider = wire.NewSet(consumer.NewWorker, LoadConsumerConfig)
)

func InitApp() (*App, func(), error) {
	panic(wire.Build(
		logProvider,
		traceProvider,
		dbProvider,
		redisProvider,
		daoProvider,
		bizProvider,
		restProvider,
		grpcProvider,
		consumerProvider, NewApp))
}
