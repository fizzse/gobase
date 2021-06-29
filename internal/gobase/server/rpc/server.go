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

func New(cfg *Config, bizCtx *biz.SampleBiz) (*Server, func(), error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, nil, err
	}

	entity := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			bizCtx.GrpcLogger(bizCtx.Logger()),
		)))

	pbBasev1.RegisterGobaseServer(entity, bizCtx)

	server := &Server{srv: entity, Listen: lis}
	return server, server.Stop, err
}

type Server struct {
	srv    *grpc.Server
	Listen net.Listener
}

func (s *Server) Run() error {
	return s.srv.Serve(s.Listen)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
