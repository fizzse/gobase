package rpc

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

	"github.com/fizzse/gobase/internal/gobase/biz"
	pbBasev1 "github.com/fizzse/gobase/protoc/gobase/v1"

	"google.golang.org/grpc"
)

/*
 * grpc server
 */

type Config struct {
	Host       string `json:"host" yaml:"host"`
	Port       int    `json:"port" yaml:"port"`
	DebugModel bool   `json:"debugModel" yaml:"debugModel"`
}

func New(cfg *Config, bizCtx *biz.SampleBiz) (instance *Server, clean func(), err error) {
	instance = &Server{bizCtx: bizCtx}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return
	}

	entity := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			instance.GrpcLogger(bizCtx.Logger()),
		)))

	pbBasev1.RegisterGobaseServer(entity, bizCtx)
	instance.Listen = lis
	instance.srv = entity
	clean = instance.Stop
	return
}

type Server struct {
	srv    *grpc.Server
	Listen net.Listener

	bizCtx *biz.SampleBiz
}

func (s *Server) register() (err error) {
	return
}

func (s *Server) Run() (err error) {
	err = s.register()
	if err != nil {
		return
	}

	return s.srv.Serve(s.Listen)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
