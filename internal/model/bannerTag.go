package model

type BannerTag struct {
	BannerID uint32 `sql:"primary_key" alias:"banner_id"`
	TagID    uint32 `sql:"primary_key" alias:"tag_id"`
}
