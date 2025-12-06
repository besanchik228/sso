package server

import (
	"context"
	pb "sso/pkg/api/test"
)

func (s *sso_server) login_handler(ctx context.Context, r *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	return s.service.Login(ctx, r)
}

func (s *sso_server) register_handler(ctx context.Context, r *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	return s.service.Register(ctx, r)
}

func (s *sso_server) refresh_handler(ctx context.Context, r *pb.RefreshTokenRequest) (*pb.RegisterUserResponse, error) {
	return s.service.Refresh(ctx, r)
}

func (s *sso_server) logout_handler(ctx context.Context, r *pb.LogOutRequest) (*pb.LogOutResponse, error) {
	return s.service.LogOut(ctx, r)
}
