package handler

import (
	"avito/internal/apperror"
	"avito/internal/dto"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserUsecase interface {
	Register(req *dto.UserRegisterRequest) error
	Login(req *dto.UserLoginRequest) (string, error)
}

type User struct {
	usecase UserUsecase
}

func NewUserHandler(usecase UserUsecase) *User {
	return &User{
		usecase: usecase,
	}
}

func (h *User) Register(ctx *gin.Context) {
	var req dto.UserRegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(fmt.Errorf("invalid request body [user handler ~ Register]: %w", apperror.ErrBadRequest))
		return
	}

	if err := h.usecase.Register(&req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *User) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(fmt.Errorf("invalid request body [user handler ~ Login]: %w", apperror.ErrBadRequest))
		return
	}

	token, err := h.usecase.Login(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": token})
}
