package rpc

import (
	"fmt"
	"net"

	"github.com/fizzse/gobase/internal/gobase/biz"
	pbBasev1 "github.com/fizzse/gobase/protoc/v1"

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

	entity := grpc.NewServer()
	pbBasev1.RegisterGobaseServer(entity, bizCtx)

	server := &Server{Entity: entity, Listen: lis}
	return server, server.Stop, err
}

type Server struct {
	Entity *grpc.Server
	Listen net.Listener
}

func (s *Server) Run() error {
	return s.Entity.Serve(s.Listen)
}

func (s *Server) Stop() {
	s.Entity.GracefulStop()
}
