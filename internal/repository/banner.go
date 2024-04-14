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

type Banner struct {
	db *sql.DB
}

func NewBannerRepository(db *sql.DB) *Banner {
	return &Banner{db: db}
}

func (r *Banner) GetByTagAndFeature(tagID, featureID int) (*model.Banner, error) {
	stmt := j.SELECT(
		t.Banner.ID,
		t.Banner.Name,
		t.Banner.IsActive,
	).FROM(
		t.Banner.
			INNER_JOIN(t.BannerTag, t.Banner.ID.EQ(t.BannerTag.BannerID)),
	).WHERE(
		t.BannerTag.TagID.EQ(j.Int(int64(tagID))).
			AND(t.Banner.FeatureID.EQ(j.Int(int64(featureID)))),
	)

	var b model.Banner
	err := stmt.Query(r.db, &b)

	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("no rows in result set [banner repository ~ GetByTagAndFeature]: %w", apperror.ErrNotFound)
		}
		return nil, fmt.Errorf("undefined error [banner repository ~ GetByTagAndFeature]: %w", apperror.ErrInternalServer)
	}

	return &b, nil
}

func (r *Banner) GetManyByTagOrFeature(tagID, featureID, limit, offset int) ([]model.Banner, error) {
	Sub := j.SELECT(
		t.Banner.ID,
	).DISTINCT(
		t.Banner.ID,
	).FROM(
		t.Banner.
			INNER_JOIN(t.BannerTag, t.Banner.ID.EQ(t.BannerTag.BannerID)),
	).WHERE(
		(t.Banner.FeatureID.EQ(j.Int(int64(featureID))).OR(
			j.Int(int64(featureID)).EQ(j.Int(0)),
		)).AND(
			t.BannerTag.TagID.EQ(j.Int(int64(tagID))).OR(
				j.Int(int64(tagID)).EQ(j.Int(0))),
		),
	).LIMIT(
		int64(limit),
	).OFFSET(
		int64(offset),
	).AsTable("sub")

	SubID := t.Banner.ID.From(Sub)

	stmt := j.SELECT(
		t.Banner.AllColumns,
		t.Tag.ID,
	).FROM(
		Sub.
			INNER_JOIN(t.Banner, SubID.EQ(t.Banner.ID)).
			INNER_JOIN(t.BannerTag, t.Banner.ID.EQ(t.BannerTag.BannerID)).
			INNER_JOIN(t.Tag, t.BannerTag.TagID.EQ(t.Tag.ID)),
	)

	var bs []model.Banner
	err := stmt.Query(r.db, &bs)

	if err != nil {
		return nil, fmt.Errorf("undefined error [banner repository ~ GetManyByTagOrFeature]: %w", apperror.ErrInternalServer)
	}
	if len(bs) == 0 {
		return nil, fmt.Errorf("no rows in result set [banner repository ~ GetManyByTagOrFeature]: %w", apperror.ErrNotFound)
	}

	return bs, nil
}

func (r *Banner) Create(b *model.Banner, tagIDs []uint32) (int, error) {
	tx, _ := r.db.Begin()

	stmt := t.Banner.INSERT(
		t.Banner.MutableColumns,
	).MODEL(
		b,
	).RETURNING(
		t.Banner.ID,
	)

	err := stmt.Query(tx, b)

	if err != nil || b.ID == 0 {
		tx.Rollback()
		return -1, fmt.Errorf("undefined error while inserting banner [banner repository ~ Create]: %w", apperror.ErrInternalServer)
	}

	bts := make([]model.BannerTag, len(tagIDs))
	for i, id := range tagIDs {
		bts[i].BannerID = b.ID
		bts[i].TagID = id
	}

	stmt = t.BannerTag.INSERT(
		t.BannerTag.AllColumns,
	).MODELS(
		bts,
	)

	_, err = stmt.Exec(tx)

	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("undefined error while inserting banner tags [banner repository ~ Create]: %w", apperror.ErrInternalServer)
	}

	tx.Commit()
	return int(b.ID), nil
}

func (r *Banner) Update(b *model.Banner, bts []model.BannerTag) error {
	tx, _ := r.db.Begin()

	exStmt := j.SELECT(
		j.EXISTS(
			j.SELECT(j.Int(1)).FROM(
				t.Banner,
			).WHERE(
				t.Banner.ID.EQ(j.Int(int64(b.ID))),
			),
		).AS("exist"),
	)

	var res struct {
		Exist bool `alias:"exist"`
	}
	err := exStmt.Query(tx, &res)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("undefined error while checking banner existance [banner repository ~ Update]: %w", apperror.ErrInternalServer)
	}
	if !res.Exist {
		tx.Rollback()
		return fmt.Errorf("banner not found [banner repository ~ Update]: %w", apperror.ErrNotFound)
	}

	updStmt := t.Banner.UPDATE(
		t.Banner.MutableColumns,
	).MODEL(
		b,
	).WHERE(
		t.Banner.ID.EQ(j.Int(int64(b.ID))),
	)

	_, err = updStmt.Exec(tx)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("undefined error while updating banner [banner repository ~ Update]: %w", apperror.ErrInternalServer)
	}

	delStmt := t.BannerTag.DELETE().WHERE(
		t.BannerTag.BannerID.EQ(j.Int(int64(b.ID))),
	)

	_, err = delStmt.Exec(tx)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("undefined error while removing banner tags [banner repository ~ Update]: %w", apperror.ErrInternalServer)
	}

	insStmt := t.BannerTag.INSERT(
		t.BannerTag.AllColumns,
	).MODELS(
		bts,
	)

	_, err = insStmt.Exec(tx)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("undefined error while inserting banner tags [banner repository ~ Update]: %w", apperror.ErrInternalServer)
	}

	tx.Commit()
	return nil
}

func (r *Banner) DeleteByID(id int) error {
	tx, _ := r.db.Begin()

	exStmt := j.SELECT(
		j.EXISTS(
			j.SELECT(j.Int(1)).FROM(
				t.Banner,
			).WHERE(
				t.Banner.ID.EQ(j.Int(int64(id))),
			),
		).AS("exist"),
	)

	var res struct {
		Exist bool `alias:"exist"`
	}
	err := exStmt.Query(tx, &res)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("undefined error while checking banner existance [banner repository ~ DeleteByID]: %w", apperror.ErrInternalServer)
	}
	if !res.Exist {
		tx.Rollback()
		return fmt.Errorf("banner not found [banner repository ~ DeleteByID]: %w", apperror.ErrNotFound)
	}

	delStmt := t.Banner.DELETE().WHERE(
		t.Banner.ID.EQ(j.Int(int64(id))),
	)

	_, err = delStmt.Exec(tx)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("undefined error [banner repository ~ DeleteByID]: %w", apperror.ErrInternalServer)
	}

	tx.Commit()
	return nil
}
