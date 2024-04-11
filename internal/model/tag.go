package model

type Tag struct {
	ID uint32 `sql:"primary_key" alias:"id"`

	Banners []Banner `alias:"banner"`
}
