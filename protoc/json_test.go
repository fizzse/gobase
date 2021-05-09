package protoc

import (
	"log"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/fizzse/gobase/protoc/gopkg"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestPb_JsonMarshal(t *testing.T) {
	a := &gopkg.SayHelloReq{Name: "fizzse", PageSize: proto.Int32(1)}
	res, _ := protojson.Marshal(a)
	log.Println(string(res))
}
