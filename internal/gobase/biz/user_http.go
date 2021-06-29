package biz

import (
	"net/http"

	pbBasev1 "github.com/fizzse/gobase/protoc/gobase/v1"

	"github.com/gin-gonic/gin"
)

func (b *SampleBiz) PingGin(c *gin.Context) {
	h := b.newRestHandler(c)
	defer h.done()

	h.Data = "pong..."
}

func (b *SampleBiz) CreateUserGin(c *gin.Context) {
	h := b.newRestHandler(c)
	defer h.done()

	req := &pbBasev1.CreateUserReq{}
	h.err = c.Bind(req)
	h.req = req

	if h.err != nil {
		h.Code = http.StatusBadRequest
		return
	}

	h.Data, h.err = b.CreateUser(c.Request.Context(), req)
	return
}
