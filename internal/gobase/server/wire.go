// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package server

import (
	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/fizzse/gobase/internal/gobase/dao"
	"github.com/fizzse/gobase/internal/gobase/server/rest"
	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
	"github.com/google/wire"
)

var (
	daoProvider  = wire.NewSet(dao.New, db.NewConn, redis.NewClient, LoadDbConfig, LoadRedisConfig)
	bizProvider  = wire.NewSet(biz.New)
	restProvider = wire.NewSet(rest.New, LoadRestConfig)
	//appProvider  = wire.NewSet()
)

func InitApp() (*App, func(), error) {
	//panic(wire.Build(Provider, dao.Provider, biz.Provider, rest.Provider, NewApp))
	panic(wire.Build(daoProvider, bizProvider, restProvider, NewApp))
}
