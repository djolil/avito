package middleware

import (
	"avito/internal/apperror"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(ctx *gin.Context) {
	ctx.Next()

	for _, e := range ctx.Errors {
		err := e.Err

		if errors.Is(err, apperror.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": apperror.ErrNotFound.Error(), "trace": err.Error()})
		} else if errors.Is(err, apperror.ErrBadRequest) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": apperror.ErrBadRequest.Error(), "trace": err.Error()})
		} else if errors.Is(err, apperror.ErrUnauthorized) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": apperror.ErrUnauthorized.Error(), "trace": err.Error()})
		} else if errors.Is(err, apperror.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": apperror.ErrForbidden.Error(), "trace": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternalServer.Error(), "trace": err.Error()})
		}
	}
}
