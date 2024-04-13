package middleware

import (
	"avito/internal/apperror"
	"fmt"

	"github.com/gin-gonic/gin"
)

func IsAdmin(ctx *gin.Context) {
	isAdmin := ctx.GetBool("is_admin")
	if !isAdmin {
		ctx.Error(fmt.Errorf("only 'admin' is allowed to perform this action [middleware ~ IsAdmin]: %w", apperror.ErrForbidden))
		ctx.Abort()
		return
	}
	ctx.Next()
}
