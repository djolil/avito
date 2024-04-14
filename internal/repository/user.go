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
		t.Role.Name,
	).FROM(
		t.UserAccount.
			LEFT_JOIN(t.UserRole, t.UserAccount.ID.EQ(t.UserRole.UserID)).
			LEFT_JOIN(t.Role, t.UserRole.RoleID.EQ(t.Role.ID)),
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
	tx, _ := r.db.Begin()

	exStmt := j.SELECT(
		j.EXISTS(
			j.SELECT(j.Int(1)).FROM(
				t.UserAccount,
			).WHERE(
				t.UserAccount.Email.EQ(j.String(u.Email)),
			),
		).AS("exist"),
	)

	var res struct {
		Exist bool `alias:"exist"`
	}
	err := exStmt.Query(tx, &res)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("undefined error while checking user existance [user repository ~ Add]: %w", apperror.ErrInternalServer)
	}
	if res.Exist {
		tx.Rollback()
		return fmt.Errorf("email already exists [user repository ~ Add]: %w", apperror.ErrBadRequest)
	}

	insStmt := t.UserAccount.INSERT(
		t.UserAccount.FirstName,
		t.UserAccount.LastName,
		t.UserAccount.Email,
		t.UserAccount.PhoneNumber,
		t.UserAccount.Password,
	).MODEL(
		u,
	)

	_, err = insStmt.Exec(tx)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("undefined error [user repository ~ Add]: %w", apperror.ErrInternalServer)
	}

	tx.Commit()
	return nil
}
