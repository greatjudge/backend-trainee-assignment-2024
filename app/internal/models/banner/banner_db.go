package banner

import (
	"encoding/json"
	"time"
)

type BannerDB struct {
	ID        int       `db:"id"`
	TagIDs    []int     `db:"tag_ids"`
	FeatureID int       `db:"feature_id"`
	Content   []byte    `db:"content"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (bDB BannerDB) ToBanner() (Banner, error) {
	b := Banner{
		ID:        bDB.ID,
		TagIDs:    bDB.TagIDs,
		FeatureID: bDB.FeatureID,
		IsActive:  bDB.IsActive,
		CreatedAt: bDB.CreatedAt,
		UpdatedAt: bDB.UpdatedAt,
	}

	err := json.Unmarshal(bDB.Content, &b.Content)
	if err != nil {
		return Banner{}, err
	}

	return b, nil
}

func SliceBannerDBToBanners(bannersDB []BannerDB) ([]Banner, error) {
	result := make([]Banner, len(bannersDB))

	for i, bDB := range bannersDB {
		b, err := bDB.ToBanner()
		if err != nil {
			return nil, err
		}

		result[i] = b
	}

	return result, nil
}
