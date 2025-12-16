package security

import (
	"context"
	"fmt"
	"sso/internal/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginLimiter struct {
	Client *redis.Client
	Limit  int
	Window time.Duration
}

func NewLoginLimiter(addr string, password string, db int, limit int, window time.Duration) *LoginLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &LoginLimiter{
		Client: rdb,
		Limit:  limit,
		Window: window,
	}
}

func (l *LoginLimiter) Check(ctx context.Context, login string) error {
	key := fmt.Sprintf("login_attempts:%s", login)
	count, err := l.Client.Incr(ctx, key).Result()
	if err != nil {
		logger.Logger().Error("error in security.Check()", zap.Error(err))
		return status.Errorf(codes.Internal, "internal")
	}
	if count == 1 {
		l.Client.Expire(ctx, key, l.Window)
	}
	if int(count) > l.Limit {
		return status.Errorf(codes.ResourceExhausted, "too many requests")
	}
	return nil
}
