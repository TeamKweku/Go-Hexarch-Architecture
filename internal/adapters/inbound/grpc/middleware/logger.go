package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/teamkweku/code-odessey-hex-arch/pkg/logger"
)

func GrpcLogger(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// create a logger with request context
		logger = logger.WithContext(ctx)

		startTime := time.Now()
		// forwards request to handler to be processed
		result, err := handler(ctx, req)
		duration := time.Since(startTime)

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		fields := map[string]interface{}{
			"protocol":    "grpc",
			"method":      info.FullMethod,
			"status_code": int(statusCode),
			"status_text": statusCode.String(),
			"duration_ms": duration.Milliseconds(),
		}

		if err != nil {
			logger.Error(ctx, err, "grpc request failed", fields)
		} else {
			logger.Info(ctx, "gRPC request proceed", fields)
		}

		return result, err
	}
}
