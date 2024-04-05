package banner

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

type BannerPartialUpdate struct {
	TagIDs    interface{} `json:"tag_ids"`
	FeatureID interface{} `json:"feature_id"`
	Content   interface{} `json:"content"`
	IsActive  interface{} `json:"is_active"`
}
