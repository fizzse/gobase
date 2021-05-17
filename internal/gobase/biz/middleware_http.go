package biz

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (b *SampleBiz) newRestHandler(ginCtx *gin.Context) *restHandler {
	handler := &restHandler{
		ginCtx:     ginCtx,
		startTime:  time.Now(),
		logger:     b.logger,
		statusCode: http.StatusOK,
	}

	handler.pre()
	return handler
}

// RestReply http 导出
type RestReply struct {
	Code      int         `json:"code"`
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
	h.Msg = "success"
	h.Timestamp = h.endTime.Unix()

	if h.err != nil {
		if h.statusCode == http.StatusOK {
			h.statusCode = http.StatusServiceUnavailable
		}

		// FIXME code 与 msg 的设定
		h.Code = http.StatusServiceUnavailable
		h.Msg = h.err.Error()
		h.Msg = "parse failed"
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
		"req", h.req,
		"res", h.Data,
		"startTime", h.startTime,
		"endTime", h.endTime,
		"duration", duration.Milliseconds(),
		"error", h.err,
	)
}
