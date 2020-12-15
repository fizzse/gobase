package rest

import (
	"fmt"
	"net/http"

	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/fizzse/gobase/pkg/loader"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var Provider = wire.NewSet(New, loader.LoadRestConfig)

func initRouter(bizCtx biz.Biz) *gin.Engine {
	route := gin.Default()
	v1 := route.Group("/gobase/v1")
	{
		v1.GET("/ping", bizCtx.Ping)
		v1.POST("/users", bizCtx.CreateUserGin)
	}

	return route
}

func New(cfg *loader.RestConfig, bizCtx biz.Biz) (*http.Server, error) {
	route := initRouter(bizCtx)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: route,
	}

	return srv, nil
}
