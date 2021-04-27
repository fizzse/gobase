// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package server

import (
	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/fizzse/gobase/internal/gobase/dao"
	"github.com/fizzse/gobase/internal/gobase/server/consumer"
	"github.com/fizzse/gobase/internal/gobase/server/rest"
	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
	"github.com/fizzse/gobase/pkg/logger"
	"github.com/google/wire"
)

var (
	logProvider      = wire.NewSet(logger.New, LoadLoggerConfig)
	dbProvider       = wire.NewSet(db.NewConn, LoadDbConfig)
	redisProvider    = wire.NewSet(redis.NewClient, LoadRedisConfig)
	daoProvider      = wire.NewSet(dao.New)
	bizProvider      = wire.NewSet(biz.New)
	restProvider     = wire.NewSet(rest.New, LoadRestConfig)
	consumerProvider = wire.NewSet(consumer.NewWorker, LoadConsumerConfig)
)

func InitApp() (*App, func(), error) {
	panic(wire.Build(
		logProvider,
		dbProvider,
		redisProvider,
		daoProvider,
		bizProvider,
		restProvider,
		consumerProvider, NewApp))
}
