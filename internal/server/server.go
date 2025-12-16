package server

import (
	service "sso/internal/service"
	pb "sso/pkg/api/test"
)

type sso_server struct {
	pb.UnimplementedSsoServer
	service *service.Service
}

func NewSsoServer(host string, port int, user, password, dbname, sslmode, redis_host, redis_password string, limit int) *sso_server {
	return &sso_server{service: service.NewService(host, port, user, password, dbname, sslmode, redis_host, redis_password, limit)}
}
