package middleware

import (
	"avito/internal/apperror"
	"avito/internal/service"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const BearerSchema = "Bearer "

type Authenticator interface {
	GetClaims(tokenString string) (*service.CustomClaims, error)
}

type AuthJWT struct {
	authcator Authenticator
}

func NewAuthJWT(authcator Authenticator) *AuthJWT {
	return &AuthJWT{
		authcator: authcator,
	}
}

func (a *AuthJWT) JWTAuth(ctx *gin.Context) {
	header := ctx.GetHeader("token")

	if header == "" {
		ctx.Error(fmt.Errorf("token header required [middleware ~ JWTAuth]: %w", apperror.ErrUnauthorized))
		ctx.Abort()
		return
	}
	if !strings.HasPrefix(header, BearerSchema) {
		ctx.Error(fmt.Errorf("invalid token header [middleware ~ JWTAuth]: %w", apperror.ErrUnauthorized))
		ctx.Abort()
		return
	}

	tokenString := strings.TrimPrefix(header, BearerSchema)

	claims, err := a.authcator.GetClaims(tokenString)
	if err != nil {
		ctx.Error(fmt.Errorf("failed to get claims [middleware ~ JWTAuth]: %w", err))
		ctx.Abort()
		return
	}

	ctx.Set("user_id", claims.UserID)
	ctx.Set("is_admin", claims.IsAdmin)
	ctx.Next()
}
