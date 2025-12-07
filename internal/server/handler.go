package server

import (
	"context"
	pb "sso/pkg/api/test"
)

func (s *sso_server) LoginUser (ctx context.Context, r *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	return s.service.Login(ctx, r)
}

func (s *sso_server) RegisterUser (ctx context.Context, r *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	return s.service.Register(ctx, r)
}

func (s *sso_server) RefreshToken (ctx context.Context, r *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	return s.service.Refresh(ctx, r)
}

func (s *sso_server) LogOut (ctx context.Context, r *pb.LogOutRequest) (*pb.LogOutResponse, error) {
	return s.service.LogOut(ctx, r)
}
