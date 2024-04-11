package model

type Banner struct {
	ID        uint32 `sql:"primary_key" alias:"id"`
	FeatureID uint32 `alias:"feature_id"`
	Name      string `alias:"name"`
	IsActive  bool   `alias:"is_active"`

	Tags []Tag `alias:"tag"`
}
