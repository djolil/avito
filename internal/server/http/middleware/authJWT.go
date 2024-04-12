package middleware

import (
	"avito/internal/apperror"
	"avito/internal/auth"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(ctx *gin.Context) {
	const BearerSchema = "Bearer "
	header := ctx.GetHeader("Authorization")

	if header == "" {
		ctx.Error(fmt.Errorf("authorization header required [middleware ~ JWTAuth]: %w", apperror.ErrUnauthorized))
		ctx.Abort()
		return
	}
	if !strings.HasPrefix(header, BearerSchema) {
		ctx.Error(fmt.Errorf("invalid authorization header [middleware ~ JWTAuth]: %w", apperror.ErrUnauthorized))
		ctx.Abort()
		return
	}

	tokenString := strings.TrimPrefix(header, BearerSchema)
	claims := auth.CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return auth.JwtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.Error(fmt.Errorf("invalid token [middleware ~ JWTAuth]: %w", apperror.ErrUnauthorized))
		ctx.Abort()
		return
	}

	ctx.Set("user_id", claims.UserID)
	ctx.Next()
}
