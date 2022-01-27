package rest

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/chenjiandongx/ginprom"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/gin-gonic/gin"
)

func (s *Server) initRouter(bizCtx *biz.SampleBiz) *gin.Engine {
	route := gin.Default()
	r := route.Group(s.cfg.Name)
	v1 := r.Group("/v1")
	{
		v1.Use(ginprom.PromMiddleware(nil))

		v1.GET("/ping", s.Ping)
		v1.POST("/users", s.CreateUser)
		v1.POST("/error", s.MockError)
	}

	// metrics
	v1.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	// profile
	ginpprof.WrapGroup(v1)
	return route
}
