package interceptor

import (
    "context"

	"sso/internal/logger"
    "go.uber.org/zap"
    "google.golang.org/grpc"
)

func LoggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    logger := logger.Logger()
    logger.Info("Incoming request",
        zap.String("method", info.FullMethod),
        zap.Any("request", req),
    )
    resp, err := handler(ctx, req)
    if err != nil {
        logger.Info("Request failed",
            zap.String("method", info.FullMethod),
            zap.Error(err),
        )
    } else {
        logger.Info("Request succeeded",
            zap.String("method", info.FullMethod),
            zap.Any("response", resp),
        )
    }

    return resp, err
}
