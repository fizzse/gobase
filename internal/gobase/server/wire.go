//go:build wireinject
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
	logProvider      = wire.NewSet(logger.New, loadLoggerConfig)
	traceProvider    = wire.NewSet(trace.New, loadTraceConfig)
	dbProvider       = wire.NewSet(db.NewConn, loadDbConfig)
	redisProvider    = wire.NewSet(redis.NewClient, loadRedisConfig)
	daoProvider      = wire.NewSet(dao.NewInstance)
	bizProvider      = wire.NewSet(biz.NewInstance)
	grpcProvider     = wire.NewSet(rpc.New, loadGrpcConfig)
	restProvider     = wire.NewSet(rest.New, loadRestConfig)
	consumerProvider = wire.NewSet(consumer.NewScheduler, loadConsumerConfig)
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
