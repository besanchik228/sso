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
    // init logger
    logger.Init()

    // Загружаем конфиг
    cfg, err := config.LoadConfig("./config/config.yaml")
    if err != nil {
        logger.Logger().Fatal("cannot load config", zap.Error(err))
    }

    // Подключаемся к БД
    if err != nil {
        logger.Logger().Fatal("cannot connect to database", zap.Error(err))
    }

    ssoSrv := server.NewSSOServer(cfg.Database.Host,
        fmt.Sprintf("%d", cfg.Database.Port),
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.DbName,
        cfg.Database.SslMode,)

    // gRPC сервер с interceptor
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

    // HTTP gateway
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
