package main

import (
    "context"
    "fmt"
    "net"
    "net/http"

	"go.uber.org/zap"  
    "google.golang.org/grpc"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "sso/internal/server"
    "sso/internal/config"
    "sso/internal/logger"
    "sso/internal/interceptor"
    pb "sso/pkg/api/test"
)

func main() {
    logger.Init()

    cfg, _ := config.LoadConfig("./config/config.yaml")

    ssoSrv := server.NewSsoServer(cfg.Database.Host,
        int(cfg.Database.Port),
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.DbName,
        cfg.Database.SslMode,
        cfg.Database.Host,
        cfg.Redis.Password,
        cfg.Redis.Limit,
    )

    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GrpcPort))
    if err != nil {
        logger.Logger().Fatal("failed to listen", zap.Error(err))
    }
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(interceptor.LoggingInterceptor),
    )
    pb.RegisterSsoServer(grpcServer, ssoSrv)

    go func() {
        logger.Logger().Sugar().Infof("gRPC server running on :%d", cfg.Server.GrpcPort)
        if err := grpcServer.Serve(lis); err != nil {
            logger.Logger().Fatal("failed to serve gRPC", zap.Error(err))
        }
    }()

    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithInsecure()}
    err = pb.RegisterSsoHandlerFromEndpoint(
        context.Background(),
        mux,
        fmt.Sprintf("localhost:%d", cfg.Server.GrpcPort),
        opts,
    )
    if err != nil {
        logger.Logger().Fatal("failed to start gateway", zap.Error(err))
    }

    logger.Logger().Sugar().Infof("HTTP gateway running on :%d", cfg.Server.HttpPort)
    if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.HttpPort), mux); err != nil {
        logger.Logger().Fatal("failed to serve HTTP", zap.Error(err))
    }
}
