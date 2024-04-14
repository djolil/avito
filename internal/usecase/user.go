package usecase

import (
	"avito/internal/apperror"
	"avito/internal/auth"
	"avito/internal/dto"
	"avito/internal/model"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByEmail(email string) (*model.UserAccount, error)
	Add(u *model.UserAccount) error
}

type JWTGenerator interface {
	GenerateJWT(userID uint32, roles []string) (string, error)
}

type User struct {
	userRepo     UserRepository
	jwtGenerator JWTGenerator
}

func NewUserUsecase(ur UserRepository, jg JWTGenerator) *User {
	return &User{
		userRepo:     ur,
		jwtGenerator: jg,
	}
}

func (u *User) Register(req *dto.UserRegisterRequest) error {
	hashPswd, err := auth.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password [user usecase ~ Register]: %w", err)
	}

	user := model.UserAccount{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashPswd,
	}

	if err := u.userRepo.Add(&user); err != nil {
		return fmt.Errorf("failed to add user [user usecase ~ Register]: %w", err)
	}
	return nil
}

func (u *User) Login(req *dto.UserLoginRequest) (string, error) {
	user, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user [user usecase ~ Login]: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("wrong password [user usecase ~ Login]: %w", apperror.ErrUnauthorized)
	}

	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
	}

	token, err := u.jwtGenerator.GenerateJWT(user.ID, roles)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT [user usecase ~ Login]: %w", apperror.ErrInternalServer)
	}

	return token, nil
}
