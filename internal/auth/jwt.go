package auth

import (
	"avito/internal/apperror"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID uint32 `json:"user_id"`
	jwt.RegisteredClaims
}

var JwtKey = []byte(os.Getenv("SECRET"))

func GenerateJWT(userID uint32) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString(JwtKey)

	if err != nil {
		return "", fmt.Errorf("failed to encode token [auth ~ GenerateJWT] %w", apperror.ErrInternalServer)
	}
	return encodedToken, nil
}
