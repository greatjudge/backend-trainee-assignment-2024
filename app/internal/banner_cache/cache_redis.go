package cache

import (
	bannermodels "banner/internal/models/banner"
	"banner/internal/service"
	"encoding/json"
	"errors"
	"time"

	"context"

	"github.com/redis/go-redis/v9"
)

// TODO change add logic
type BannerRedisCache struct {
	client           *redis.Client
	bannerExpiration time.Duration
}

func NewBannerRedisCahe(client *redis.Client, bannerExpiration time.Duration) *BannerRedisCache {
	return &BannerRedisCache{
		client:           client,
		bannerExpiration: bannerExpiration,
	}
}

func (c *BannerRedisCache) GetBanner(ctx context.Context, tagID int, featureID int) (bannermodels.Banner, error) {
	data, err := c.client.Get(ctx, formKeyFromTagIDFeatureID(tagID, featureID)).Result()

	switch {
	case errors.Is(err, redis.Nil):
		return bannermodels.Banner{}, service.ErrCacheBannerNotFound
	case err != nil:
		return bannermodels.Banner{}, err
	}

	var banner bannermodels.Banner
	err = json.Unmarshal([]byte(data), &banner)
	if err != nil {
		return bannermodels.Banner{}, err
	}

	return banner, nil
}

func (c *BannerRedisCache) SetBanner(ctx context.Context, tagID int, featureID int, banner bannermodels.Banner) error {
	bannerBytes, err := json.Marshal(banner)
	if err != nil {
		return err
	}

	err = c.client.Set(
		ctx,
		formKeyFromTagIDFeatureID(tagID, featureID),
		bannerBytes,
		c.bannerExpiration,
	).Err()

	return err
}
