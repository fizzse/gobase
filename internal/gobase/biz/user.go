package biz

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID      uint
	Account string
}

func (b *SampleBiz) CreateUserGin(ginCtx *gin.Context) {
	req := &CreateUserReq{}
	if err := ginCtx.Bind(req); err != nil {
		ginCtx.Status(http.StatusBadRequest)
		return
	}

	res, err := b.CreateUser(ginCtx.Request.Context(), req)
	if err != nil {

	}

	ginCtx.JSON(http.StatusOK, res)
}

func (b *SampleBiz) CreateUser(ctx context.Context, user *CreateUserReq) (*UserInfo, error) {
	// b.daoCtx.CreateUser()
	return nil, nil
}
