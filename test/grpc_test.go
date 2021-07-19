package test

import (
	"context"
	"log"
	"testing"

	"google.golang.org/protobuf/proto"

	pbBasev1 "github.com/fizzse/gobase/protoc/gobase/v1"
	"google.golang.org/grpc"
)

func TestGrpcCli(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal("dial failed", err)
	}

	cli := pbBasev1.NewGobaseClient(conn)

	res, err := cli.CreateUser(context.Background(), &pbBasev1.CreateUserReq{Name: "fizzse", Age: proto.Int64(27)})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
}
