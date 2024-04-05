package service

import (
	bannermodels "banner/internal/models/banner"
	usermodels "banner/internal/models/user"
	"context"
)

type bannerRepo interface {
	GetBanner(ctx context.Context, tagID int, featureID int, useLastVersion bool) (bannermodels.Banner, error)
	GetFiltered(ctx context.Context, filter bannermodels.FilterSchema) ([]bannermodels.Banner, error)
	CreateBanner(ctx context.Context, banner bannermodels.Banner) (int, error)
	PartialUpdateBanner(ctx context.Context, bannerPartial bannermodels.BannerPartialUpdate) error
	DeleteBanner(ctx context.Context, id int) error
}

type BannerService struct {
	repo bannerRepo
}

func NewBannerService(bannerRepo bannerRepo) *BannerService {
	return &BannerService{
		repo: bannerRepo,
	}
}

func (s BannerService) GetUserBanner(ctx context.Context, user usermodels.User, tagID int, featureID int, useLastRevision bool) (map[string]interface{}, error) {
	b, err := s.repo.GetBanner(ctx, tagID, featureID, useLastRevision) // Разделить запрос банера и контента?
	if err != nil {
		return nil, err // TODO
	}

	if !b.IsActive && !user.IsAdmin {
		return nil, ErrBannerNotFound
	}

	return b.Content, nil
}

func (s BannerService) BannerList(ctx context.Context, filter bannermodels.FilterSchema) ([]bannermodels.Banner, error) {
	banners, err := s.repo.GetFiltered(ctx, filter)
	if err != nil {
		return nil, err // TODO
	}
	return banners, nil
}

func (s BannerService) CreateBanner(ctx context.Context, banner bannermodels.Banner) (int, error) {
	id, err := s.repo.CreateBanner(ctx, banner)
	if err != nil {
		return 0, nil // TODO
	}
	return id, nil
}

func (s BannerService) PartialUpdateBanner(ctx context.Context, id int, bannerPartial bannermodels.BannerPartialUpdate) error {
	err := s.repo.PartialUpdateBanner(ctx, bannerPartial)
	if err != nil {
		return err
	}
	return nil
}

func (s BannerService) DeleteBanner(ctx context.Context, id int) error {
	err := s.repo.DeleteBanner(ctx, id)
	if err != nil {
		return err // TODO
	}
	return nil
}
