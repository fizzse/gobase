package rest

import (
	"fmt"
	"net/http"

	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/gin-gonic/gin"
)

//var Provider = wire.NewSet(New, loader.LoadRestConfig)

type Config struct {
	Host       string `json:"host" yaml:"host"`
	Port       int    `json:"port" yaml:"port"`
	DebugModel bool   `json:"debugModel" yaml:"debugModel"`
}

func initRouter(bizCtx biz.Biz) *gin.Engine {
	route := gin.Default()
	v1 := route.Group("/gobase/v1")
	{
		v1.GET("/ping", bizCtx.Ping)
		v1.POST("/users", bizCtx.CreateUserGin)
	}

	return route
}

func New(cfg *Config, bizCtx biz.Biz) (*http.Server, error) {
	if !cfg.DebugModel {
		gin.SetMode(gin.ReleaseMode)
	}

	route := initRouter(bizCtx)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: route,
	}

	return srv, nil
}
