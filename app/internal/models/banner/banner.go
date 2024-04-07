package banner

import (
	"time"
)

type Banner struct {
	ID        int                    `json:"banner_id"`
	TagIDs    []int                  `json:"tag_ids"`
	FeatureID int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func UpdatedBanner(banner Banner, bannerPartial BannerPartialUpdate) (Banner, error) {
	if bannerPartial.FeatureID != nil {
		featureID, ok := bannerPartial.FeatureID.(int)
		if !ok {
			return Banner{}, ErrBadFeatureID
		}
		banner.FeatureID = featureID
	}

	if bannerPartial.TagIDs != nil {
		tagIDs, ok := bannerPartial.TagIDs.([]int)
		if !ok {
			return Banner{}, ErrBadTagIDs
		}
		banner.TagIDs = tagIDs
	}

	if bannerPartial.Content != nil {
		content, ok := bannerPartial.Content.(map[string]interface{})
		if !ok {
			return Banner{}, ErrBadContent
		}
		banner.Content = content
	}

	if bannerPartial.IsActive != nil {
		isActive, ok := bannerPartial.IsActive.(bool)
		if !ok {
			return Banner{}, ErrBadIsActive
		}
		banner.IsActive = isActive
	}

	return banner, nil
}
