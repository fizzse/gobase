package biz

import (
	"net/http"

	pbBasev1 "github.com/fizzse/gobase/protoc/v1"

	"github.com/gin-gonic/gin"
)

func (b *SampleBiz) PingGin(c *gin.Context) {
	c.String(http.StatusOK, "pong...")
}

func (b *SampleBiz) CreateUserGin(c *gin.Context) {
	req := &pbBasev1.CreateUserReq{}
	if err := c.Bind(req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	res, err := b.CreateUser(c.Request.Context(), req)
	if err != nil {

	}

	c.JSON(http.StatusOK, res)
}
