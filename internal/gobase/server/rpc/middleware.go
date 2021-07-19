package rpc

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (s *Server) GrpcLogger(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logging := logger.Infow

		startTime := time.Now()
		span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
		defer span.Finish()

		resp, err := handler(ctx, req)
		duration := time.Since(startTime)
		endTime := time.Now()
		if err != nil {
			logging = logger.Errorw // 修改日志级别为错误
		}

		logging(info.FullMethod,
			"traceId", span,
			"req", req,
			"res", resp,
			"startTime", startTime,
			"endTime", endTime,
			"duration", duration.Milliseconds(),
			"error", err,
		)

		return resp, err
	}
}

func (s *Server) GrpcTrace() grpc.UnaryServerInterceptor {
	return nil
}
