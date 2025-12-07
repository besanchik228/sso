package service

import (
	"context"
	"sso/internal/auth"
	repo "sso/internal/repository"
	pb "sso/pkg/api/test"
)

type Service struct {
	repo *repo.Repository
}

func NewService(host, port, user, password, dbname, sslmode string) *Service {
	repository, err := repo.NewRepository(host, port, user, password, dbname, sslmode)
	if err != nil {
		return nil
	}
	return &Service{repo: repository}
}

func (s *Service) Login(ctx context.Context, r *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	hash_password, err := auth.HashPassword(r.Password)
	if err != nil {
		return nil, err
	}
	u, err := s.repo.CheckUser(r.Login, hash_password)
	if err != nil {
		return nil, err
	}
	access_token, refresh_token := auth.NewToken(u.Id, u.Login)
	return &pb.LoginUserResponse{
		AccessToken: access_token,
		RefreshToken: refresh_token,
		AccessTokenTime: 60 * 60 * 24,
		Err: "",
	}, nil

}

func (s *Service) Register(ctx context.Context, r *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	hash_password, err := auth.HashPassword(r.Password)
	if err != nil {
		return nil, err
	}
	u, err := s.repo.CreateUser(r.Login, hash_password)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterUserResponse{Ok: true, User: &pb.UserPublic{
		Id: u.Id,
		Login: u.Login,
	}}, nil
}

func (s *Service) Refresh(ctx context.Context, r *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	userID, login, _ := s.repo.GetUserIDByRefreshToken(r.RefreshToken)
	access_token, err := auth.Refresh(userID, login, r.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &pb.RefreshTokenResponse{
		AccessToken: access_token,
		AccessTokenTime: 15,
	}, nil
}

func (s *Service) LogOut(ctx context.Context, r *pb.LogOutRequest) (*pb.LogOutResponse, error) {
	err := s.repo.RevokeRefresh(r.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &pb.LogOutResponse{Ok: true}, nil
}
