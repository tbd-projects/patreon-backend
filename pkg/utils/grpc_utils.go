package utils

import (
	"context"
	"patreon/pkg/monitoring"
	"time"

	"google.golang.org/grpc"
)

func AuthInterceptor(metrics monitoring.Monitoring) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		reply, err := handler(ctx, req)

		statusCode := "200"
		if err != nil {
			statusCode = "500"
			metrics.GetErrorsHits().WithLabelValues(statusCode, info.FullMethod, info.FullMethod).Inc()
		} else {
			metrics.GetSuccessHits().WithLabelValues(statusCode, info.FullMethod, info.FullMethod).Inc()
		}

		metrics.GetExecution().
			WithLabelValues(statusCode, info.FullMethod, info.FullMethod).
			Observe(time.Since(start).Seconds())
		metrics.GetRequestCounter().Inc()

		return reply, err

	}
}
