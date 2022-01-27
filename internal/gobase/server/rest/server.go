package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	DebugModel bool   `yaml:"debugModel"`
	Name       string `yaml:"name"`
}

type Server struct {
	cfg    *Config
	srv    *http.Server
	bizCtx *biz.SampleBiz
}

func New(cfg *Config, bizCtx *biz.SampleBiz) (instance *Server, err error) {
	if !cfg.DebugModel {
		gin.SetMode(gin.ReleaseMode)
	}

	instance = &Server{cfg: cfg, bizCtx: bizCtx}

	route := instance.initRouter(bizCtx)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: route,
	}

	instance.srv = srv
	return
}

func (s *Server) Name() string {
	return "rest"
}

func (s *Server) Run(ctx context.Context) (err error) {
	err = s.register()
	if err != nil {
		return
	}

	err = s.srv.ListenAndServe()
	return err
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // TODO timeout config
	defer cancel()

	_ = s.srv.Shutdown(ctx)
	return
}

func (s *Server) register() (err error) {
	return
}
