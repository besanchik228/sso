package service

import (
	"context"
	repo "sso/internal/repository"
	pb "sso/pkg/api/test"
)

type Service struct {
	repo *repo.Repository
}

func NewService() *Service {
	return &Service{repo: repo.NewRepository()}
}

func (s *Service) Login(ctx context.Context, r *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

}

func (s *Service) Register(ctx context.Context, r *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {

}

func (s *Service) Refresh(ctx context.Context, r *pb.RefreshTokenRequest) (*pb.RegisterUserResponse, error) {

}

func (s *Service) LogOut(ctx context.Context, r *pb.LogOutRequest) (*pb.LogOutResponse, error) {

}
