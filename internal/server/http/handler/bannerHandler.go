package handler

import (
	"avito/internal/apperror"
	"avito/internal/dto"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BannerUsecase interface {
	GetByTagAndFeature(tagID, featureID int) (*dto.BannerResponse, error)
	GetManyByTagOrFeature(tagID, featureID, limit, offset int) ([]dto.BannerDetailsResponse, error)
	Create(req *dto.BannerCreateRequest) (*dto.BannerCreateResponse, error)
	Update(id int, req *dto.BannerUpdateRequest) error
	Delete(id int) error
}

type Banner struct {
	usecase BannerUsecase
}

func NewBannerHandler(usecase BannerUsecase) *Banner {
	return &Banner{
		usecase: usecase,
	}
}

func (h *Banner) GetByTagAndFeature(ctx *gin.Context) {
	tagID, err := strconv.Atoi(ctx.Query("tag_id"))
	if err != nil {
		ctx.Error(fmt.Errorf("invalid 'tag_id' parameter [banner handler ~ GetByTagAndFeature]: %w", apperror.ErrBadRequest))
		return
	}
	featureID, err := strconv.Atoi(ctx.Query("feature_id"))
	if err != nil {
		ctx.Error(fmt.Errorf("invalid 'feature_id' parameter [banner handler ~ GetByTagAndFeature]: %w", apperror.ErrBadRequest))
		return
	}

	b, err := h.usecase.GetByTagAndFeature(tagID, featureID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, b)
}

func (h *Banner) GetManyByTagOrFeature(ctx *gin.Context) {
	tagID, err := strconv.Atoi(ctx.Query("tag_id"))
	if err != nil {
		ctx.Error(fmt.Errorf("invalid 'tag_id' parameter [banner handler ~ GetManyByTagOrFeature]: %w", apperror.ErrBadRequest))
		return
	}
	featureID, err := strconv.Atoi(ctx.Query("feature_id"))
	if err != nil {
		ctx.Error(fmt.Errorf("invalid 'feature_id' parameter [banner handler ~ GetManyByTagOrFeature]: %w", apperror.ErrBadRequest))
		return
	}
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		ctx.Error(fmt.Errorf("invalid 'limit' parameter [banner handler ~ GetManyByTagOrFeature]: %w", apperror.ErrBadRequest))
		return
	}
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		ctx.Error(fmt.Errorf("invalid 'offset' parameter [banner handler ~ GetManyByTagOrFeature]: %w", apperror.ErrBadRequest))
		return
	}

	bs, err := h.usecase.GetManyByTagOrFeature(tagID, featureID, limit, offset)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, bs)
}

func (h *Banner) Create(ctx *gin.Context) {
	var req dto.BannerCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(fmt.Errorf("invalid request body [banner handler ~ Create]: %w", apperror.ErrBadRequest))
		return
	}

	b, err := h.usecase.Create(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, b)
}

func (h *Banner) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(fmt.Errorf("invalid banner 'id' parameter [banner handler ~ Update]: %w", apperror.ErrBadRequest))
		return
	}

	var req dto.BannerUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(fmt.Errorf("invalid request body [banner handler ~ Update]: %w", apperror.ErrBadRequest))
		return
	}

	if err := h.usecase.Update(id, &req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Banner) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(fmt.Errorf("invalid banner 'id' parameter [banner handler ~ Delete]: %w", apperror.ErrBadRequest))
		return
	}

	if err := h.usecase.Delete(id); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
