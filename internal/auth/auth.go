package auth

import (
	"time"

	"sso/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte("")

func Init(s string) {
	secret = []byte(s)
}

func NewToken(userID, login string) (string, string, bool) {
	ok := true
	accessData := jwt.MapClaims{
		"sub":   userID,
		"login": login,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessData)
	accessToken, err := access.SignedString(secret)
	if err != nil {
		logger.Logger().Error("access_token creation error (from auth.NewToken())")
		ok = false
	}
	refreshData := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshData)
	refreshToken, err := refresh.SignedString(secret)
	if err != nil {
		logger.Logger().Error("refresh_token creation error (from auth.NewToken())")
		ok = false
	}
	
	return accessToken, refreshToken, ok
}

func Refresh(userID, login, refreshToken string) (string, error, bool) {
	ok := true
	accessData := jwt.MapClaims{
		"sub":   userID,
		"login": login,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessData)
	accessToken, err := access.SignedString(secret)
	if err != nil {
		logger.Logger().Error("access_token creation error (from auth.Refresh())")
		ok = false
	}
	return accessToken, nil, ok
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 31)
	return string(bytes), err
}

func CheckPassword(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
