package biz

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (b *SampleBiz) GrpcLogger(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logging := logger.Infow

		startTime := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(startTime)
		endTime := time.Now()
		if err != nil {
			logging = logger.Errorw // 修改日志级别为错误
		}

		logging(info.FullMethod,
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

func (b *SampleBiz) GrpcTrace() grpc.UnaryServerInterceptor {
	return nil
}
