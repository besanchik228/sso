package server

import (
	service "sso/internal/service"
	pb "sso/pkg/api/test"
)

type sso_server struct {
	pb.UnimplementedSsoServer
	service *service.Service
}

func NewSSOServer(host, port, user, password, dbname, sslmode string) *sso_server {
	return &sso_server{service: service.NewService(host, port, user, password, dbname, sslmode)}
}
