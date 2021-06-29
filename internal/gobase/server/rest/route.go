package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DeanThompson/ginpprof"
	"github.com/chenjiandongx/ginprom"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Host       string `json:"host" yaml:"host"`
	Port       int    `json:"port" yaml:"port"`
	DebugModel bool   `json:"debugModel" yaml:"debugModel"`
}

type Server struct {
	cfg *Config
	srv *http.Server
}

func New(cfg *Config, bizCtx *biz.SampleBiz) (*Server, error) {
	if !cfg.DebugModel {
		gin.SetMode(gin.ReleaseMode)
	}

	route := initRouter(bizCtx)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: route,
	}

	server := &Server{cfg: cfg, srv: srv}
	return server, nil
}

func (s *Server) Run() (err error) {
	err = s.srv.ListenAndServe()
	return err
}

func (s *Server) Stop(ctx context.Context) (err error) {
	err = s.srv.Shutdown(ctx)
	return
}

func initRouter(bizCtx *biz.SampleBiz) *gin.Engine {
	route := gin.Default()
	v1 := route.Group("/gobase/v1")
	{
		v1.Use(ginprom.PromMiddleware(nil))

		v1.GET("/ping", bizCtx.PingGin)
		v1.POST("/users", bizCtx.CreateUserGin)
	}

	// metrics
	v1.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	// profile
	ginpprof.WrapGroup(v1)
	return route
}
