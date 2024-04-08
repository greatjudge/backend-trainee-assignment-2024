package cache

import (
	bannermodels "banner/internal/models/banner"
	"banner/internal/service"
	"context"
)

// TODO change add logic
type BannerCahe struct {
}

func NewBannerCache() *BannerCahe {
	return &BannerCahe{}
}

func (c *BannerCahe) GetBanner(ctx context.Context, tagID int, featureID int) (bannermodels.Banner, error) {
	return bannermodels.Banner{}, service.ErrCacheBannerNotFound
}

func (c *BannerCahe) SetBanner(ctx context.Context, tagID int, featureID int, banner bannermodels.Banner) error {
	return nil
}
