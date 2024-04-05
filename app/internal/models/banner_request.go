package models

type BannerRequest struct {
	TagIDs    []int                  `json:"tag_ids"`
	FeatureID int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
}

func (br BannerRequest) ToBanner() Banner {
	return Banner{
		TagIDs:    br.TagIDs,
		FeatureID: br.FeatureID,
		Content:   br.Content,
		IsActive:  br.IsActive,
	}
}
