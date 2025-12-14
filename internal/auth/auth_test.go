package auth

import (
    "testing"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

func TestNewToken(t *testing.T) {
    access, refresh, ok := NewToken("123", "testuser")
    if !ok {
        t.Fatalf("expected ok=true, got false")
    }
    if access == "" || refresh == "" {
        t.Fatalf("expected non-empty tokens, got access=%q refresh=%q", access, refresh)
    }
    parsed, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
        return secret, nil
    })
    if err != nil || !parsed.Valid {
        t.Fatalf("access token invalid: %v", err)
    }
}

func TestRefresh(t *testing.T) {
    _, refresh, _ := NewToken("123", "testuser")
    access, err, ok := Refresh("123", "testuser", refresh)
    if !ok {
        t.Fatalf("expected ok=true, got false")
    }
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if access == "" {
        t.Fatalf("expected non-empty access token")
    }
    parsed, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
        return secret, nil
    })
    if err != nil || !parsed.Valid {
        t.Fatalf("refreshed access token invalid: %v", err)
    }
}

func TestPasswordHashing(t *testing.T) {
    password := "supersecret"
    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("unexpected error hashing password: %v", err)
    }
    if hash == "" {
        t.Fatalf("expected non-empty hash")
    }
    if err := CheckPassword(hash, password); err != nil {
        t.Fatalf("expected password to match, got error: %v", err)
    }
    if err := CheckPassword(hash, "wrongpassword"); err == nil {
        t.Fatalf("expected error for wrong password, got nil")
    }
}

func TestTokenExpiration(t *testing.T) {
    access, _, _ := NewToken("123", "testuser")

    parsed, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
        return secret, nil
    })
    if err != nil {
        t.Fatalf("parse error: %v", err)
    }
    claims, ok := parsed.Claims.(jwt.MapClaims)
    if !ok {
        t.Fatalf("expected MapClaims")
    }
    exp := int64(claims["exp"].(float64))
    if exp <= time.Now().Unix() {
        t.Fatalf("expected exp in future, got %d", exp)
    }
}
