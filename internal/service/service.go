package service

import (
	"context"
	"sso/internal/auth"
	"sso/internal/logger"
	repo "sso/internal/repository"
	pb "sso/pkg/api/test"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repo *repo.Repository
}

func NewService(host string, port int, user, password, dbname, sslmode string) *Service {
	repository, err := repo.NewRepository(host, port, user, password, dbname, sslmode)
	if err != nil {
		logger.Logger().Fatal("service creation error (from repo.NewRepository())", zap.Error(err))
		return nil
	}
	return &Service{repo: repository,}
}

func (s *Service) Login(ctx context.Context, r *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	u, err := s.repo.CheckUser(r.Login, r.Password)
	if err != nil {
		return nil, err
	}
	access_token, refresh_token, ok := auth.NewToken(u.Id, u.Login)
	if !ok {
		logger.Logger().Error("creation token error (from auth.NewToken())", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal")
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
		logger.Logger().Error("hash password error (from auth.HashPassword())", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal")
	}
	u, err := s.repo.CreateUser(r.Login, hash_password)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterUserResponse{Ok: true, Err: "", User: &pb.UserPublic{
		Id: u.Id,
		Login: u.Login,
	}}, nil
}

func (s *Service) Refresh(ctx context.Context, r *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	userID, login, err := s.repo.GetUserIDByRefreshToken(r.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "refresh token not found")
	}
	if !s.repo.ValidToken(r.RefreshToken) {
		return nil, status.Errorf(codes.InvalidArgument, "token has expired")
	}
	access_token, err, ok := auth.Refresh(userID, login, r.RefreshToken)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal")
	}
	if err != nil {
		logger.Logger().Error("refresh password error (from auth.Refresh())", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal")
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
