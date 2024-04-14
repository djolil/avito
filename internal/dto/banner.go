package dto

type BannerResponse struct {
	Name     string `json:"name"`
	IsActive bool   `json:"-"`
}

type BannerDetailsResponse struct {
	BannerID  uint32   `json:"banner_id"`
	TagIDs    []uint32 `json:"tag_ids"`
	FeatureID uint32   `json:"feature_id"`
	Content   Banner   `json:"content"`
	IsActive  bool     `json:"is_active"`
}

type Banner struct {
	Name string `json:"name" binding:"required"`
}

type BannerCreateRequest struct {
	TagIDs    []uint32 `json:"tag_ids" binding:"required"`
	FeatureID uint32   `json:"feature_id" binding:"required"`
	Content   Banner   `json:"content" binding:"required"`
	IsActive  bool     `json:"is_active"`
}

type BannerCreateResponse struct {
	BannerID uint32 `json:"banner_id"`
}

type BannerUpdateRequest struct {
	TagIDs    []uint32 `json:"tag_ids" binding:"required"`
	FeatureID uint32   `json:"feature_id" binding:"required"`
	Content   Banner   `json:"content" binding:"required"`
	IsActive  bool     `json:"is_active"`
}
