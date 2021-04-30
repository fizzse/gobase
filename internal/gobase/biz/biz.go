package biz

import (
	"context"
	"net/http"

	"github.com/fizzse/gobase/pkg/mq/kafka"
	pbBase "github.com/fizzse/gobase/protoc/gopkg"

	"github.com/fizzse/gobase/internal/gobase/dao"
	"github.com/gin-gonic/gin"
)

// 显示声明 SampleBiz 实现了Biz
var _ Biz = &SampleBiz{}

// GinBiz http interface
type GinBiz interface {
	Ping(ginCtx *gin.Context)
	CreateUserGin(ginCtx *gin.Context)
}

type Biz interface {
	GinBiz              // http server
	pbBase.GobaseServer // grpc

	Close()
	CreateUser(ctx context.Context, user *CreateUserReq) (*UserInfo, error)

	DealMsg(ctx context.Context, msg kafka.Event) error
}

// New 返回抽象的接口
func New(daoCtx dao.Dao) (Biz, func(), error) {
	bizCtx := &SampleBiz{
		daoCtx: daoCtx,
	}

	return bizCtx, bizCtx.Close, nil
}

type SampleBiz struct {
	daoCtx                           dao.Dao
	pbBase.UnimplementedGobaseServer // FIXME 默认实现所有api
}

func (b *SampleBiz) Close() {
	b.daoCtx.Close()
}

func (b *SampleBiz) Ping(ginCtx *gin.Context) {
	ginCtx.String(http.StatusOK, "pong...")
}

func (b *SampleBiz) SendMsgToDevice(ctx context.Context, in *pbBase.SayHelloReq) (*pbBase.SayHelloRes, error) {
	return &pbBase.SayHelloRes{Reply: "hello" + in.Name}, nil
}
