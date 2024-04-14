package service

import (
	"avito/internal/apperror"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID  uint32 `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

type JWTAuth struct {
	jwtKey []byte
}

type ConfigJWT struct {
	SecretKey string `yaml:"secret_key"`
}

func NewJWTAuth(cfg *ConfigJWT) *JWTAuth {
	return &JWTAuth{
		jwtKey: []byte(cfg.SecretKey),
	}
}

func (j *JWTAuth) GenerateJWT(userID uint32, roles []string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}
	for _, r := range roles {
		switch r {
		case "admin":
			claims.IsAdmin = true
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString(j.jwtKey)

	if err != nil {
		return "", fmt.Errorf("failed to encode token [jwtauth service ~ GenerateJWT] %w", apperror.ErrInternalServer)
	}
	return encodedToken, nil
}

func (j *JWTAuth) GetClaims(tokenString string) (*CustomClaims, error) {
	claims := CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return j.jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token [jwtauth service ~ GetClaims]: %w", apperror.ErrUnauthorized)
	}

	return &claims, nil
}
