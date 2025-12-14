package repository

import (
	"database/sql"
	"fmt"
	"sso/internal/auth"
	"sso/internal/logger"
	"sso/internal/models"
	"time"

    "go.uber.org/zap"
	_ "github.com/jackc/pgx/v4/stdlib"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Repository struct {
    db *sql.DB
}

func NewRepository(host string, port int, user, password, dbname, sslmode string) (*Repository, error) {
    dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
        user, password, host, port, dbname, sslmode)
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        logger.Logger().Fatal("error repository.NewRepository()", zap.Error(err))
        return nil, err
    }
    if err := db.Ping(); err != nil {
        logger.Logger().Fatal("error repository.NewRepository()", zap.Error(err))
        return nil, err
    }
    return &Repository{db: db}, nil
}

func (r *Repository) CreateUser(login, passwordHash string) (*models.User, error) {
    var id string
    err := r.db.QueryRow(
        `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`,
        login, passwordHash,
    ).Scan(&id)
    if err != nil {
        return nil, status.Errorf(codes.AlreadyExists, "user with that login already exists")
    }
    return &models.User{
        Id:    id,
        Login: login,
        PasswordHash:  passwordHash,
    }, nil
}

func (r *Repository) CheckUser(login, password string) (*models.User, error) {
    u := &models.User{}
    err := r.db.QueryRow(
        `SELECT id, login, password_hash FROM users WHERE login=$1`,
        login,
    ).Scan(&u.Id, &u.Login, &u.PasswordHash)
    
    if err != nil {
        if err != sql.ErrNoRows {
            logger.Logger().Error("error in repository.CheckUser", zap.Error(err))
            return nil, status.Errorf(codes.Internal, "internal")
        }
        return nil, status.Errorf(codes.Unauthenticated, "login not found")
    }

    if err = auth.CheckPassword(u.PasswordHash, password); err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "wrong password")
    }

    return u, nil
}

func (r *Repository) RevokeRefresh(token string) error {
    exist := false
    err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM refresh_tokens WHERE token=$1)", token).Scan(&exist)
    if !exist || err == sql.ErrNoRows {
        return status.Errorf(codes.InvalidArgument, "token not found")
    }
    if err != nil {
        logger.Logger().Error("error in repository.RevokeRefresh()", zap.Error(err))
    }
    _, err = r.db.Exec(`UPDATE refresh_tokens SET revoked=true WHERE token=$1`, token)
    return err
}

func (r *Repository) NewRefreshToken(token string, userID string, expiresAt time.Time) error {
    _, err := r.db.Exec(
        `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
        userID, token, expiresAt,
    )
    return err
}

func (r *Repository) GetUserIDByRefreshToken(token string) (string, string, error){
	u := &models.User{}
	err := r.db.QueryRow(
		"SELECT user_id FROM refresh_tokens WHERE token=$1", token,
	).Scan(&u.Id)
	if err != nil {
        if err != sql.ErrNoRows{
            logger.Logger().Error("error in repository.GetUserIDByRefreshToken", zap.Error(err))
        }
		return "", "", err
	}
	err = r.db.QueryRow(
		"SELECT login FROM users WHERE id=$1", u.Id, 
	).Scan(&u.Login)
	if err != nil {
        if err != sql.ErrNoRows{
            logger.Logger().Error("error in repository.GetUserIDByRefreshToken", zap.Error(err))
        }
		return "", "", err
	}
	return u.Id, u.Login, nil
}

func (r *Repository) ValidToken(token string) bool {
    var exp time.Time
    var revoked bool
    err := r.db.QueryRow(
        "SELECT expires_at, revoked FROM refresh_tokens WHERE token=$1",
        token,
    ).Scan(&exp, &revoked)

    if err != nil {
        if err == sql.ErrNoRows {
            return false
        }
        logger.Logger().Error("error in repository.ValidToken()", zap.Error(err))
        return false
    }
    if time.Now().UTC().Before(exp) && !revoked {
        return true
    }
    return false
}
