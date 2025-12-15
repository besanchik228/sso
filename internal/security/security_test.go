package security_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sso/internal/security"
)

func newTestLimiter(addr string, limit int, window time.Duration) *security.LoginLimiter {
    rdb := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    return &security.LoginLimiter{
		Client: rdb,
		Limit: limit,
		Window: window,
    }
}

func TestLoginLimiter_AllowsWithinLimit(t *testing.T) {
    s, err := miniredis.Run()
    if err != nil {
        t.Fatalf("failed to start miniredis: %v", err)
    }
    defer s.Close()

    limiter := newTestLimiter(s.Addr(), 3, time.Minute)
    ctx := context.Background()

    for i := 0; i < 3; i++ {
        if err := limiter.Check(ctx, "user"); err != nil {
            t.Errorf("unexpected error on attempt %d: %v", i+1, err)
        }
    }
}

func TestLoginLimiter_BlocksAfterLimit(t *testing.T) {
    s, err := miniredis.Run()
    if err != nil {
        t.Fatalf("failed to start miniredis: %v", err)
    }
    defer s.Close()

    limiter := newTestLimiter(s.Addr(), 2, time.Minute)
    ctx := context.Background()

    if err := limiter.Check(ctx, "user"); err != nil {
        t.Fatalf("unexpected error on first attempt: %v", err)
    }
    if err := limiter.Check(ctx, "user"); err != nil {
        t.Fatalf("unexpected error on second attempt: %v", err)
    }

    // третья должна вернуть ResourceExhausted
    err = limiter.Check(ctx, "user")
    if err == nil {
        t.Fatal("expected error, got nil")
    }
    st, _ := status.FromError(err)
    if st.Code() != codes.ResourceExhausted {
        t.Errorf("expected ResourceExhausted, got %v", st.Code())
    }
}

func TestLoginLimiter_SetsTTL(t *testing.T) {
    s, err := miniredis.Run()
    if err != nil {
        t.Fatalf("failed to start miniredis: %v", err)
    }
    defer s.Close()

    limiter := newTestLimiter(s.Addr(), 1, 30*time.Second)
    ctx := context.Background()

    if err := limiter.Check(ctx, "user"); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // проверяем, что ключ имеет TTL
    ttl := s.TTL("login_attempts:user")
    if ttl <= 0 {
        t.Errorf("expected TTL > 0, got %v", ttl)
    }
}
