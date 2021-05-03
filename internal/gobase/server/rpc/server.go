package rpc

import (
	"fmt"
	"net"

	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/fizzse/gobase/protoc/gopkg"

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

func New(cfg *Config, bizCtx *biz.SampleBiz) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}

	entity := grpc.NewServer()
	gopkg.RegisterGobaseServer(entity, bizCtx)
	return &Server{Entity: entity, Listen: lis}, err
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
