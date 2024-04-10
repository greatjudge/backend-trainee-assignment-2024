package cache

import (
	bannermodels "banner/internal/models/banner"
	"banner/internal/service"
	"context"
)

// cache that do nothing, if want to no caching
type BannerNoCache struct {
}

func NewBannerNoCache() *BannerNoCache {
	return &BannerNoCache{}
}

func (c *BannerNoCache) GetBanner(ctx context.Context, tagID int, featureID int) (bannermodels.Banner, error) {
	return bannermodels.Banner{}, service.ErrCacheBannerNotFound
}

func (c *BannerNoCache) SetBanner(ctx context.Context, tagID int, featureID int, banner bannermodels.Banner) error {
	return nil
}
