package service

import (
	bannermodels "banner/internal/models/banner"
	usermodels "banner/internal/models/user"
	"context"
	"encoding/json"
	"errors"
)

type bannerRepo interface {
	GetUserBanner(ctx context.Context, tagID int, featureID int) (bannermodels.Banner, error)
	GetFiltered(ctx context.Context, filter bannermodels.FilterSchema) ([]bannermodels.Banner, error)
	CreateBanner(ctx context.Context, banner bannermodels.Banner) (int, error)
	PartialUpdateBanner(ctx context.Context, id int, bannerPartial bannermodels.BannerPartialUpdate) error
	DeleteBanner(ctx context.Context, id int) error
}

type bannerCache interface {
	GetBanner(ctx context.Context, tagID int, featureID int) (bannermodels.Banner, error)
	SetBanner(ctx context.Context, tagID int, featureID int, banner bannermodels.Banner) error
}

type BannerService struct {
	repo  bannerRepo
	cache bannerCache
}

func NewBannerService(bannerRepo bannerRepo, bannerCache bannerCache) *BannerService {
	return &BannerService{
		repo:  bannerRepo,
		cache: bannerCache,
	}
}

// get banner from cahe and if not exists get from repo and set to cache
func (s *BannerService) getOrSetUserBannerFromCache(ctx context.Context, tagID int, featureID int) (bannermodels.Banner, error) {
	b, err := s.cache.GetBanner(ctx, tagID, featureID)

	// NO err
	if err == nil {
		return b, nil
	}

	if !errors.Is(err, ErrCacheBannerNotFound) {
		return bannermodels.Banner{}, err
	}

	// If not found in cache
	b, err = s.repo.GetUserBanner(ctx, tagID, featureID)

	switch {
	case errors.Is(err, ErrDBBannerNotFound):
		return bannermodels.Banner{}, ErrBannerNotFound
	case err != nil:
		return bannermodels.Banner{}, err
	}

	err = s.cache.SetBanner(ctx, tagID, featureID, b)
	if err != nil {
		return bannermodels.Banner{}, err
	}

	return b, nil
}

func (s *BannerService) GetUserBanner(ctx context.Context, user usermodels.User, tagID int, featureID int, useLastRevision bool) ([]byte, error) {
	var b bannermodels.Banner
	var err error

	if useLastRevision {
		b, err = s.repo.GetUserBanner(ctx, tagID, featureID)

		switch {
		case errors.Is(err, ErrDBBannerNotFound):
			return nil, ErrBannerNotFound
		case err != nil:
			return nil, err
		}
	} else {
		b, err = s.getOrSetUserBannerFromCache(ctx, tagID, featureID)

		if err != nil {
			return nil, err
		}
	}

	if !b.IsActive && !user.IsAdmin {
		return nil, ErrBannerNotFound
	}

	contentJSON, err := json.Marshal(b.Content)
	if err != nil {
		return nil, err // TODO
	}

	return contentJSON, nil
}

func (s *BannerService) BannerList(ctx context.Context, filter bannermodels.FilterSchema) ([]bannermodels.Banner, error) {
	banners, err := s.repo.GetFiltered(ctx, filter)
	if err != nil {
		return nil, err
	}
	return banners, nil
}

func (s *BannerService) CreateBanner(ctx context.Context, banner bannermodels.Banner) (int, error) {
	id, err := s.repo.CreateBanner(ctx, banner)

	switch {
	case errors.Is(err, ErrDBBannerAlreadyExists):
		return 0, ErrBannerAlreadyExists
	case err != nil:
		return 0, err
	}

	return id, nil
}

func (s *BannerService) PartialUpdateBanner(ctx context.Context, id int, bannerPartial bannermodels.BannerPartialUpdate) error {
	err := s.repo.PartialUpdateBanner(ctx, id, bannerPartial)

	switch {
	case errors.Is(err, ErrDBBannerAlreadyExists):
		return ErrBannerAlreadyExists
	case err != nil:
		return err
	}

	return nil
}

func (s *BannerService) DeleteBanner(ctx context.Context, id int) error {
	err := s.repo.DeleteBanner(ctx, id)

	switch {
	case errors.Is(err, ErrDBBannerNotFound):
		return ErrBannerNotFound
	case err != nil:
		return err
	}

	return nil
}
