package biz

import (
	"context"
	"net/http"

	"github.com/fizzse/gobase/internal/gobase/dao"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var Provider = wire.NewSet(New)

type GinBiz interface {
	Ping(ginCtx *gin.Context)
	CreateUserGin(ginCtx *gin.Context)
}

//
type Biz interface {
	GinBiz // http server
	Close()
	CreateUser(ctx context.Context, user *CreateUserReq) (*UserInfo, error)
}

func New(daoCtx dao.Dao) (Biz, func(), error) {
	bizCtx := &SampleBiz{
		daoCtx: daoCtx,
	}

	return bizCtx, bizCtx.Close, nil
}

type SampleBiz struct {
	daoCtx dao.Dao
}

func (b *SampleBiz) Close() {
	b.daoCtx.Close()
}

func (b *SampleBiz) Ping(ginCtx *gin.Context) {
	ginCtx.String(http.StatusOK, "pong...")
}
