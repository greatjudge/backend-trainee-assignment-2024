package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	bannermodels "banner/internal/models/banner"
	"banner/internal/service"
)

const (
	bannerTableName = "banner"

	stmtCreateBannerTemplate = `
	INSERT INTO %v (tag_ids, feature_id, content, is_active)
	VALUES ($1, $2, $3, $4)
	WHERE NOT EXISTS (SELECT id FROM %v WHERE feature_id=$2 AND tag_ids && $1) 
	RETURNING id;
	`

	stmtGetUserBannerTemplae = `
	SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at
	FROM %v WHERE feature_id=$1 AND $2 = ANY(tag_ids);
	`
)

var (
	stmtCreateBanner  = fmt.Sprintf(stmtCreateBannerTemplate, bannerTableName, bannerTableName)
	stmtGetUserBanner = fmt.Sprintf(stmtGetUserBannerTemplae, bannerTableName)
)

type BannerRepo struct {
	db *pgxpool.Pool
}

func (repo BannerRepo) CreateBanner(ctx context.Context, banner bannermodels.Banner) (int, error) {
	row := repo.db.QueryRow(ctx, stmtCreateBanner, banner.TagIDs, banner.FeatureID, banner.Content, banner.IsActive)

	var id int
	err := row.Scan(&id)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return 0, service.ErrBannerAlreadyExists
	case err != nil:
		return 0, err
	}

	return id, nil
}

func (repo BannerRepo) GetUserBanner(ctx context.Context, tagID int, featureID int, useLastRevision bool) (bannermodels.Banner, error) {
	row := repo.db.QueryRow(ctx, stmtGetUserBanner, tagID, featureID)

	var contentJSON []byte
	var banner bannermodels.Banner

	err := row.Scan(
		banner.ID,
		banner.TagIDs,
		banner.FeatureID,
		contentJSON,
		banner.IsActive,
		banner.CreatedAt,
		banner.UpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return bannermodels.Banner{}, service.ErrBannerNotFound
	case err != nil:
		return bannermodels.Banner{}, err
	}

	var content map[string]interface{}
	err = json.Unmarshal(contentJSON, &content)
	if err != nil {
		return bannermodels.Banner{}, err
	}

	banner.Content = content

	return banner, nil
}

func (repo BannerRepo) PartialUpdateBanner(ctx context.Context, id int, bannerPartial bannermodels.BannerPartialUpdate) error {
	return nil
}
