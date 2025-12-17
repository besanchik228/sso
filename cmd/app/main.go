package main

import (
    "fmt"
    "net"

    "go.uber.org/zap"
    "google.golang.org/grpc"

    "sso/internal/server"
    "sso/internal/config"
    "sso/internal/logger"
    "sso/internal/interceptor"
    pb "sso/pkg/api/test"
)

func main() {
    logger.Init()

    cfg, err := config.LoadConfig("./config/config.yaml")
    if err != nil {
        logger.Logger().Fatal("failed to load config", zap.Error(err))
    }

    ssoSrv := server.NewSsoServer(
        cfg.Database.Host,
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

    logger.Logger().Sugar().Infof("SSO gRPC server running on :%d", cfg.Server.GrpcPort)
    if err := grpcServer.Serve(lis); err != nil {
        logger.Logger().Fatal("failed to serve gRPC", zap.Error(err))
    }
}
