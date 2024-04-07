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

	stmtSelectBannerTemplate = `
	SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at
	FROM %v"
	`

	stmtCheckBannerConsistentTemplate = `
	SELECT id FROM %v WHERE id != $1 AND feature_id=$2 AND $3 = ANY(tag_ids);
	`

	stmtUpdateBannerTemplate = `
	UPDATE %v SET tag_ids=$1, feature_id=$2, content=$3, is_active=$4 WHERE id=$5
	`
)

var (
	stmtSelectBanner          = fmt.Sprintf(stmtSelectBannerTemplate, bannerTableName)
	stmtCreateBanner          = fmt.Sprintf(stmtCreateBannerTemplate, bannerTableName, bannerTableName)
	stmtGetUserBanner         = stmtSelectBanner + " WHERE feature_id=$1 AND $2 = ANY(tag_ids);"
	stmtCheckBannerConsistent = fmt.Sprintf(stmtCheckBannerConsistentTemplate, bannerTableName)
	stmtGetBaner              = stmtSelectBanner + " WHERE id=$1"
	stmtUpdateBanner          = fmt.Sprintf(stmtUpdateBannerTemplate, bannerTableName)
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
		return 0, service.ErrDBBannerAlreadyExists
	case err != nil:
		return 0, err
	}

	return id, nil
}

func (repo BannerRepo) GetUserBanner(ctx context.Context, tagID int, featureID int) (bannermodels.Banner, error) {
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
		return bannermodels.Banner{}, service.ErrDBBannerNotFound
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
	return repo.partialUpdateBanner(
		ctx,
		id,
		bannerPartial,
		bannerPartial.FeatureID != nil || bannerPartial.TagIDs != nil,
	)
}

func (repo BannerRepo) partialUpdateBanner(ctx context.Context, id int, bannerPartial bannermodels.BannerPartialUpdate, checkConsistent bool) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	row := repo.db.QueryRow(ctx, stmtGetBaner, id)

	var contentJSON []byte
	var banner bannermodels.Banner

	err = row.Scan(
		banner.ID,
		banner.TagIDs,
		banner.FeatureID,
		contentJSON,
		banner.IsActive,
		banner.CreatedAt,
		banner.UpdatedAt,
	)

	if err != nil {
		tx.Rollback(ctx)

		if errors.Is(err, pgx.ErrNoRows) {
			return service.ErrBannerNotFound
		}
		return err
	}

	updatedBanner, err := bannermodels.UpdatedBanner(banner, bannerPartial)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if checkConsistent {
		row = repo.db.QueryRow(
			ctx,
			stmtCheckBannerConsistent,
			updatedBanner.ID,
			updatedBanner.FeatureID,
			updatedBanner.TagIDs,
		)

		var idBad int
		err = row.Scan(&idBad)
		// NO err
		if err == nil {
			tx.Rollback(ctx)
			return service.ErrDBBannerAlreadyExists
		}
	}

	res, err := repo.db.Exec(
		ctx,
		stmtUpdateBanner,
		updatedBanner.TagIDs,
		updatedBanner.FeatureID,
		updatedBanner.Content,
		updatedBanner.IsActive,
		updatedBanner.ID,
	)

	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if res.RowsAffected() == 0 {
		tx.Rollback(ctx)
		return fmt.Errorf("while updating banner with id=%v has deleted", id)
	}

	return nil
}
