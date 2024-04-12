package repository

import (
	"avito/internal/apperror"
	"avito/internal/model"
	"database/sql"
	"errors"
	"fmt"

	t "avito/internal/jet/table"

	j "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

type User struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *User {
	return &User{db: db}
}

func (r *User) GetByEmail(email string) (*model.UserAccount, error) {
	stmt := j.SELECT(
		t.UserAccount.AllColumns,
	).FROM(
		t.UserAccount,
	).WHERE(
		t.UserAccount.Email.EQ(j.String(email)),
	)

	var u model.UserAccount
	err := stmt.Query(r.db, &u)

	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("no rows in result set [user repository ~ GetByEmail]: %w", apperror.ErrNotFound)
		}
		return nil, fmt.Errorf("undefined error [user repository ~ GetByEmail]: %w", apperror.ErrInternalServer)
	}

	return &u, nil
}

func (r *User) Add(u *model.UserAccount) error {
	exists, err := r.existsByEmail(u.Email)

	if err != nil {
		return fmt.Errorf("failed to check if user exists [user repository ~ Add]: %w", err)
	}
	if exists {
		return fmt.Errorf("email already exists [user repository ~ Add]: %w", apperror.ErrBadRequest)
	}

	stmt := t.UserAccount.INSERT(
		t.UserAccount.FirstName,
		t.UserAccount.LastName,
		t.UserAccount.Email,
		t.UserAccount.PhoneNumber,
		t.UserAccount.Password,
	).MODEL(
		u,
	)

	_, err = stmt.Exec(r.db)

	if err != nil {
		return fmt.Errorf("undefined error [user repository ~ Add]: %w", apperror.ErrInternalServer)
	}

	return nil
}

func (r *User) existsByEmail(email string) (bool, error) {
	stmt := j.SELECT(
		j.EXISTS(
			j.SELECT(j.Int(1)).FROM(
				t.UserAccount,
			).WHERE(
				t.UserAccount.Email.EQ(j.String(email)),
			),
		).AS("exists"),
	)

	var res struct {
		Exists bool `alias:"exists"`
	}
	err := stmt.Query(r.db, &res)

	if err != nil {
		return false, fmt.Errorf("undefined error [user repository ~ existsByEmail]: %w", apperror.ErrInternalServer)
	}

	return res.Exists, nil
}
