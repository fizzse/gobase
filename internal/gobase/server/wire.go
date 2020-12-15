// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package server

import (
	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/fizzse/gobase/internal/gobase/dao"
	"github.com/fizzse/gobase/internal/gobase/server/rest"
	"github.com/google/wire"
)

func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, biz.Provider, rest.Provider, NewApp))
}
