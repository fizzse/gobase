package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	qerror "github.com/fizzse/gobase/pkg/errors"
	pbBasev1 "github.com/fizzse/gobase/protoc/gobase/v1"
)

func (s *Server) newRestHandler(ginCtx *gin.Context) *restHandler {
	handler := &restHandler{
		ginCtx:     ginCtx,
		startTime:  time.Now(),
		logger:     s.bizCtx.Logger(),
		statusCode: http.StatusOK,
	}

	handler.pre()
	return handler
}

// RestReply http 导出
type RestReply struct {
	Code      int         `json:"code"`
	Reason    string      `json:"reason"`
	Msg       string      `json:"msg"`
	TraceId   string      `json:"traceId"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// 小写是因为这些字段不需要导出
type restHandler struct {
	ginCtx *gin.Context
	logger *zap.SugaredLogger

	opName     string
	req        interface{}
	err        error
	statusCode int
	startTime  time.Time
	endTime    time.Time
	span       opentracing.Span

	RestReply
}

type ValidatorStruct interface {
	Validate() error
}

func (h *restHandler) ShouldBindJSON(req interface{}) (err error) {
	h.req = req
	err = h.ginCtx.ShouldBindJSON(req)
	if err != nil {
		err = qerror.BadRequest(pbBasev1.ERR_CODE_INPUT_DATA_FORMAT_ERROR.String(), err.Error())
		h.err = err
		return
	}

	if v, ok := req.(ValidatorStruct); ok { // 参数验证
		err = v.Validate()
		if err != nil {
			err = qerror.BadRequest(pbBasev1.ERR_CODE_INPUT_DATA_FORMAT_ERROR.String(), err.Error())
			h.err = err
			return
		}
	}

	return
}

func (h *restHandler) preTrace() {
	remoteContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(h.ginCtx.Request.Header),
	)

	opName := h.ginCtx.Request.Method + ":" + h.ginCtx.FullPath()
	h.opName = opName
	if err != nil {
		h.span = opentracing.StartSpan(opName)
	} else {
		h.span = opentracing.StartSpan(opName, opentracing.ChildOf(remoteContext))
	}

	h.ginCtx.Request = h.ginCtx.Request.WithContext(opentracing.ContextWithSpan(h.ginCtx.Request.Context(), h.span))
}

func (h *restHandler) pre() {
	h.preTrace()
}

func (h *restHandler) done() {
	h.endTime = time.Now()

	h.Code = 200
	h.Reason = "success"
	h.Msg = "success"
	h.Timestamp = h.endTime.Unix()

	if h.err != nil {
		h.statusCode = qerror.Code(h.err)

		reason := qerror.Reason(h.err)
		code := pbBasev1.ERR_CODE_value[reason]
		h.Code = int(code)
		h.Reason = reason
		h.Msg = qerror.Message(h.err)
	}

	h.TraceId = fmt.Sprintf("%v", h.span)
	h.reply()
	h.logging()
	h.span.Finish()
}

func (h *restHandler) reply() {
	h.ginCtx.JSON(h.statusCode, h.RestReply)
}

func (h *restHandler) logging() {
	logging := h.logger.Infow
	if h.err != nil {
		logging = h.logger.Errorw
	}

	duration := time.Since(h.startTime)

	logging(h.opName,
		"traceId", h.TraceId,
		"req", h.req,
		"res", h.Data,
		"startTime", h.startTime,
		"endTime", h.endTime,
		"duration", duration.Milliseconds(),
		"error", h.err,
	)
}
