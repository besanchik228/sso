package service

import (
	"sso/internal/logger"
	"context"
	"sso/internal/auth"
	repo "sso/internal/repository"
	pb "sso/pkg/api/test"
	"time"
)

type Service struct {
	repo *repo.Repository
}

func NewService(host string, port int, user, password, dbname, sslmode string) *Service {
	repository, err := repo.NewRepository(host, port, user, password, dbname, sslmode)
	if err != nil {
		logger.Logger().Fatal("service creation error")
		return nil
	}
	return &Service{repo: repository,}
}

func (s *Service) Login(ctx context.Context, r *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	u, err := s.repo.CheckUser(r.Login, r.Password)
	if err != nil {
		return &pb.LoginUserResponse{
			AccessToken: "",
			AccessTokenTime: 0,
			RefreshToken: "",
			Err: "401 error / 16 error / UNAUTHENTICATED",
		}, nil
	}
	access_token, refresh_token, ok := auth.NewToken(u.Id, u.Login)
	if !ok {
		return &pb.LoginUserResponse{
			AccessToken: "",
			AccessTokenTime: 0,
			RefreshToken: "",
			Err: "500 error / 13 error / INTERNAL",
		}, nil
	}
	s.repo.NewRefreshToken(refresh_token, u.Id, time.Now().Add(time.Hour * 24))
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
		logger.Logger().Error("hash password error (from auth.HashPassword())")
		return &pb.RegisterUserResponse{
			Ok: false,
			User: nil,
			Err: "500 error / 13 error / INTERNAL",
		}, nil
	}
	u, err := s.repo.CreateUser(r.Login, hash_password)
	if err != nil {
		return &pb.RegisterUserResponse{Ok: false, Err: "409 error / 6 error / ALREADY_EXISTS", User: nil}, err
	}
	return &pb.RegisterUserResponse{Ok: true, Err: "", User: &pb.UserPublic{
		Id: u.Id,
		Login: u.Login,
	}}, nil
}

func (s *Service) Refresh(ctx context.Context, r *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	userID, login, err := s.repo.GetUserIDByRefreshToken(r.RefreshToken)
	if err != nil {
		return &pb.RefreshTokenResponse{
			AccessToken: "",
			AccessTokenTime: 0,
			Err: "401 error / 16 error / UNAUTHENTICATED",
		}, nil
	}
	access_token, err, ok := auth.Refresh(userID, login, r.RefreshToken)
	if !ok {
		return &pb.RefreshTokenResponse{
			AccessToken: "",
			AccessTokenTime: 0,
			Err: "500 error / 13 error / INTERNAL",
		}, nil
	}
	if err != nil {
		logger.Logger().Error("refresh password error (from auth.Refresh())")
		return &pb.RefreshTokenResponse{
			AccessToken: "",
			AccessTokenTime: 0,
			Err: "500 error / 13 error / INTERNAL",
		}, nil
	}
	return &pb.RefreshTokenResponse{
		AccessToken: access_token,
		AccessTokenTime: 15,
	}, nil
}

func (s *Service) LogOut(ctx context.Context, r *pb.LogOutRequest) (*pb.LogOutResponse, error) {
	err := s.repo.RevokeRefresh(r.RefreshToken)
	if err != nil {
		return &pb.LogOutResponse{Ok: false}, nil
	}
	return &pb.LogOutResponse{Ok: true}, nil
}
