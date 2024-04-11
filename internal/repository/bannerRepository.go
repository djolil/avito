package repository

import (
	"avito/internal/model"
	"database/sql"

	t "avito/internal/jet/table"

	j "github.com/go-jet/jet/v2/postgres"
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
		return nil, err
	}
	return &b, nil
}

func (r *Banner) GetManyByTagOrFeature(tagID, featureID, limit, offset int) ([]model.Banner, error) {
	stmt := j.SELECT(
		t.Banner.AllColumns,
		t.Tag.ID,
	).FROM(
		t.Banner.
			INNER_JOIN(t.BannerTag, t.Banner.ID.EQ(t.BannerTag.BannerID)).
			INNER_JOIN(t.Tag, t.BannerTag.TagID.EQ(t.Tag.ID)),
	).WHERE(
		(t.Banner.FeatureID.EQ(j.Int(int64(featureID))).
			OR(j.Int(int64(featureID)).EQ(j.Int(0)))).
			AND(
				j.EXISTS(
					j.SELECT(j.Int(1)).FROM(
						t.BannerTag,
					).WHERE(
						t.BannerTag.BannerID.EQ(t.Banner.ID).
							AND(t.BannerTag.TagID.EQ(j.Int(int64(tagID))).
								OR(j.Int(int64(tagID)).EQ(j.Int(0)))),
					),
				),
			),
	).LIMIT(
		int64(limit),
	).OFFSET(
		int64(offset),
	)

	var bs []model.Banner
	err := stmt.Query(r.db, &bs)

	if err != nil {
		return nil, err
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
		return -1, err
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
		return -1, err
	}

	tx.Commit()
	return int(b.ID), nil
}

func (r *Banner) Update(b *model.Banner, bts []model.BannerTag) error {
	tx, _ := r.db.Begin()

	updStmt := t.Banner.UPDATE(
		t.Banner.MutableColumns,
	).MODEL(
		b,
	).WHERE(
		t.Banner.ID.EQ(j.Int(int64(b.ID))),
	)

	_, err := updStmt.Exec(tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	delStmt := t.BannerTag.DELETE().WHERE(
		t.BannerTag.BannerID.EQ(j.Int(int64(b.ID))),
	)

	_, err = delStmt.Exec(tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	insStmt := t.BannerTag.INSERT(
		t.BannerTag.AllColumns,
	).MODELS(
		bts,
	)

	_, err = insStmt.Exec(tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Banner) DeleteByID(id int) error {
	stmt := t.Banner.DELETE().WHERE(
		t.Banner.ID.EQ(j.Int(int64(id))),
	)

	_, err := stmt.Exec(r.db)

	if err != nil {
		return err
	}
	return nil
}
