package server

import (
	service "sso/internal/service"
	pb "sso/pkg/api/test"
)

type sso_server struct {
	pb.UnimplementedSsoServer
	service *service.Service
}

func NewSSOServer() *sso_server {
	return &sso_server{service: service.NewService()}
}
