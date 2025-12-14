package repository

import (
	"context"
	"sso/internal/auth"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newMockRepo(t *testing.T) (*Repository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed creation sqlmock: %v", err)
	}
	return &Repository{db: db}, mock
}

func TestCreateUser(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

    repo, mock := newMockRepo(t)

    mock.ExpectQuery(`INSERT INTO users`).
        WithArgs("login", "hash").
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))

    u, err := repo.CreateUser("login", "hash", ctx)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if u.Id != "123" || u.Login != "login" {
        t.Errorf("unexpected user: %+v", u)
    }
}

func TestCheckUserSuccess(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

    repo, mock := newMockRepo(t)

    hash, _ := auth.HashPassword("secret")

    mock.ExpectQuery(`SELECT id, login, password_hash FROM users WHERE login=`).
        WithArgs("login").
        WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).
            AddRow("123", "login", hash))

    u, err := repo.CheckUser("login", "secret", ctx)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if u.Login != "login" {
        t.Errorf("expected login=login, got %s", u.Login)
    }
}

func TestCheckUserWrongPassword(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

    repo, mock := newMockRepo(t)

    hash, _ := auth.HashPassword("secret")

    mock.ExpectQuery(`SELECT id, login, password_hash FROM users WHERE login=`).
        WithArgs("login").
        WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).
            AddRow("123", "login", hash))

    _, err := repo.CheckUser("login", "wrong", ctx)
    if status.Code(err) != codes.Unauthenticated {
        t.Errorf("expected Unauthenticated, got %v", err)
    }
}

func TestRevokeRefresh(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

    repo, mock := newMockRepo(t)

    mock.ExpectQuery(`SELECT EXISTS`).
        WithArgs("token").
        WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

    mock.ExpectExec(`UPDATE refresh_tokens SET revoked=true`).
        WithArgs("token").
        WillReturnResult(sqlmock.NewResult(1, 1))

    err := repo.RevokeRefresh("token", ctx)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestNewRefreshToken(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

    repo, mock := newMockRepo(t)

    mock.ExpectExec(`INSERT INTO refresh_tokens`).
        WithArgs("user123", "token123", sqlmock.AnyArg()).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err := repo.NewRefreshToken("token123", "user123", time.Now().Add(time.Hour), ctx)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestGetUserIDByRefreshToken(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

    repo, mock := newMockRepo(t)

    mock.ExpectQuery(`SELECT user_id FROM refresh_tokens`).
        WithArgs("token123").
        WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow("user123"))

    mock.ExpectQuery(`SELECT login FROM users WHERE id=`).
        WithArgs("user123").
        WillReturnRows(sqlmock.NewRows([]string{"login"}).AddRow("login123"))

    id, login, err := repo.GetUserIDByRefreshToken("token123", ctx)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if id != "user123" || login != "login123" {
        t.Errorf("unexpected result: id=%s login=%s", id, login)
    }
}

func TestValidToken(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

    repo, mock := newMockRepo(t)

    exp := time.Now().Add(time.Hour)

    mock.ExpectQuery(`SELECT expires_at, revoked FROM refresh_tokens`).
        WithArgs("token123").
        WillReturnRows(sqlmock.NewRows([]string{"expires_at", "revoked"}).
            AddRow(exp, false))

    if !repo.ValidToken("token123", ctx) {
        t.Errorf("expected token to be valid")
    }
}