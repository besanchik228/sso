package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("Incoming request: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	log.Printf("Response: %v, Error: %v", resp, err)
	return resp, err
}
