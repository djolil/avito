package usecase

import (
	"avito/internal/dto"
	"avito/internal/model"
	"fmt"
)

type BannerRepository interface {
	GetByTagAndFeature(tagID, featureID int) (*model.Banner, error)
	GetManyByTagOrFeature(tagID, featureID, limit, offset int) ([]model.Banner, error)
	Create(b *model.Banner, tagIDs []uint32) (int, error)
	Update(b *model.Banner, bts []model.BannerTag) error
	DeleteByID(id int) error
}

type Banner struct {
	bannerRepo BannerRepository
}

func NewBannerUsecase(br BannerRepository) *Banner {
	return &Banner{
		bannerRepo: br,
	}
}

func (u *Banner) GetByTagAndFeature(tagID, featureID int) (*dto.BannerResponse, error) {
	b, err := u.bannerRepo.GetByTagAndFeature(tagID, featureID)
	if err != nil {
		return nil, fmt.Errorf("failed to get banner [banner usecase ~ GetByTagAndFeature]: %w", err)
	}
	res := dto.BannerResponse{
		Name:     b.Name,
		IsActive: b.IsActive,
	}
	return &res, nil
}

func (u *Banner) GetManyByTagOrFeature(tagID, featureID, limit, offset int) ([]dto.BannerDetailsResponse, error) {
	bs, err := u.bannerRepo.GetManyByTagOrFeature(tagID, featureID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get banners [banner usecase ~ GetManyByTagOrFeature]: %w", err)
	}

	res := make([]dto.BannerDetailsResponse, len(bs))
	for i, b := range bs {
		cur := &res[i]
		cur.BannerID = b.ID
		cur.FeatureID = b.FeatureID
		cur.Content.Name = b.Name
		cur.IsActive = b.IsActive

		cur.TagIDs = make([]uint32, len(b.Tags))
		for i, t := range b.Tags {
			cur.TagIDs[i] = t.ID
		}
	}
	return res, nil
}

func (u *Banner) Create(req *dto.BannerCreateRequest) (*dto.BannerCreateResponse, error) {
	b := model.Banner{
		FeatureID: req.FeatureID,
		Name:      req.Content.Name,
		IsActive:  req.IsActive,
	}

	bID, err := u.bannerRepo.Create(&b, req.TagIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create banner [banner usecase ~ Create]: %w", err)
	}

	res := dto.BannerCreateResponse{
		BannerID: uint32(bID),
	}
	return &res, nil
}

func (u *Banner) Update(id int, req *dto.BannerUpdateRequest) error {
	b := model.Banner{
		ID:        uint32(id),
		FeatureID: req.FeatureID,
		Name:      req.Content.Name,
		IsActive:  req.IsActive,
	}

	bts := make([]model.BannerTag, len(req.TagIDs))
	for i, tagID := range req.TagIDs {
		bts[i].BannerID = uint32(id)
		bts[i].TagID = tagID
	}

	if err := u.bannerRepo.Update(&b, bts); err != nil {
		return fmt.Errorf("failed to update banner [banner usecase ~ Update]: %w", err)
	}
	return nil
}

func (u *Banner) Delete(id int) error {
	if err := u.bannerRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("failed to delete banner [banner usecase ~ Delete]: %w", err)
	}
	return nil
}
