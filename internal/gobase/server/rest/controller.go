package rest

import (
	"net/http"

	pbBasev1 "github.com/fizzse/gobase/protoc/gobase/v1"

	"github.com/gin-gonic/gin"
)

/*
 http controller
 You should rename this file according to your business
 example: user.go
*/

func (s *Server) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong...")
}

func (s *Server) CreateUser(c *gin.Context) {
	h := s.newRestHandler(c)
	defer h.done()

	req := &pbBasev1.CreateUserReq{}
	h.err = c.Bind(req)
	h.req = req

	if h.err != nil {
		h.Code = http.StatusBadRequest
		return
	}

	h.Data, h.err = s.bizCtx.CreateUser(c.Request.Context(), req)
	return
}
