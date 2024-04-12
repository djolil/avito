package auth

import (
	"avito/internal/apperror"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("failed to get hash for password: %w", apperror.ErrInternalServer)
	}
	return string(bytes), nil
}
