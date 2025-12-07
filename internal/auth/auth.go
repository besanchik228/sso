package auth

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("secret")

func NewToken(userID, login string) (string, string) {
	accessData := jwt.MapClaims{
		"sub" : userID,
		"login" : login,
		"exp" : time.Now().Add(time.Minute * 15).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodES256, accessData)
	accessToken, _ := access.SignedString(secret)
	refreshData := jwt.MapClaims{
		"sub" : userID,
		"exp" : time.Now().Add(time.Hour * 24).Unix(),
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodES256, refreshData)
	refreshToken, _ := refresh.SignedString(secret)
	return accessToken, refreshToken
}

func Refresh(userID, login, refreshToken string) (string, error){
	accessData := jwt.MapClaims{
		"sub" : userID,
		"login" : login,
		"exp" : time.Now().Add(time.Minute * 15).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodES256, accessData)
	accessToken, _ := access.SignedString(secret)
	return accessToken, nil
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}