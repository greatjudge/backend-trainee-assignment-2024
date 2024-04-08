package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	bannermodels "banner/internal/models/banner"
	"banner/internal/service"
	"banner/internal/tools"
)

const (
	SQLDuplicateErrCode = "23505"

	stmtCreateBanner = `
	with create_banner AS (
		INSERT into banner (is_active, "content") VALUES ($3, $4) RETURNING "id"
	),
	create_banner_relation as (
		INSERT into banner_relation (banner_id, feature_id, tag_id)
		SELECT create_banner.id as banner_id, $2 as feature_id, UNNEST($1) as tag_id FROM create_banner
	)
	  
	SELECT "id" FROM create_banner;
	`

	stmtGetUserBanner = `
	with find_banner as (
		SELECT banner_id, tag_id, feature_id FROM banner_relation WHERE feature_id=$2 AND tag_id=$1
	)
	
	SELECT
		b.id,
		fb.feature_id,
		(SELECT ARRAY_AGG(tag_id) FROM banner_relation AS br WHERE br.banner_id = b.id) as tag_ids,
		b.is_active,
		b.created_at,
		b.updated_at
	FROM banner as b JOIN find_banner as fb ON (b.id = fb.banner_id);
	`

	stmtGetBanerByID = `
	SELECT
		b.id,
		fb.feature_id,
		(SELECT ARRAY_AGG(tag_id) FROM banner_relation AS br WHERE br.banner_id = b.id) as tag_ids,
		b.is_active,
		b.created_at,
		b.updated_at
	FROM banner as b
	WHERE b.id = $1;
	`

	stmtUpdateBanner = `
	UPDATE banner SET is_active=$2, "content"=$3 WHERE "id"=$1;
	`

	stmtDeleteOldTagIDs = `
	DELETE from banner_relation WHERE banner_id=$1 AND feature_id=$2 AND tag_id = ANY($3);
	`

	stmtInsertNewTagIDs = `
	INSERT INTO banner_relation (banner_id, feature_id, tag_id) SELECT $1, $2, UNNEST($3);
	`

	stmtUpdateFeatureID = `
	UPDATE banner_relation SET feature_id=$2 WHERE banner_id=$1;
	`

	stmtBannerList = `
	SELECT
		b.id,
		fb.feature_id,
		(SELECT ARRAY_AGG(tag_id) FROM banner_relation AS br WHERE br.banner_id = b.id) as tag_ids,
		b.is_active,
		b.created_at,
		b.updated_at
	FROM banner as b
	ORDER BY b.created_at DESC
	LIMIT $1 OFFSET $2
	`

	stmtBannerListWithFilterTemplate = `
	with filtered_banners as (
	  SELECT banner_id from banner_relation WHERE %v
	)
	
	SELECT
		b.id,
		fb.feature_id,
		(SELECT ARRAY_AGG(tag_id) FROM banner_relation AS br WHERE br.banner_id = b.id) as tag_ids,
		b.is_active,
		b.created_at,
		b.updated_at
	FROM banner as b JOIN filtered_banners as fb ON (b.id = fb.banner_id)
	order by b.created_at DESC
	LIMIT $1 OFFSET $2;
	`
)

type BannerRepo struct {
	db *pgxpool.Pool
}

func (repo BannerRepo) CreateBanner(ctx context.Context, banner bannermodels.Banner) (int, error) {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	contentJSON, err := json.Marshal(banner.Content)
	if err != nil {
		return 0, err
	}

	row := tx.QueryRow(
		ctx,
		stmtCreateBanner,
		banner.TagIDs,
		banner.FeatureID,
		contentJSON,
		banner.IsActive,
	)

	var id int
	err = row.Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == SQLDuplicateErrCode {
			return 0, service.ErrBannerAlreadyExists
		}

		return 0, err
	}

	tx.Commit(ctx)
	return id, nil
}

func (repo BannerRepo) GetUserBanner(ctx context.Context, tagID int, featureID int) (bannermodels.Banner, error) {
	// ADD TRANSATION ?
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
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, stmtGetBanerByID, id)

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
		if errors.Is(err, pgx.ErrNoRows) {
			return service.ErrBannerNotFound
		}
		return err
	}

	err = json.Unmarshal(contentJSON, &banner.Content)
	if err != nil {
		return err
	}

	updatedBanner, err := bannermodels.UpdatedBanner(banner, bannerPartial)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}

	if bannerPartial.TagIDs != nil {
		// banner.TagIDs is old tags
		toDelete := tools.SliceDiff(banner.TagIDs, updatedBanner.TagIDs)
		toInsert := tools.SliceDiff(updatedBanner.TagIDs, banner.TagIDs)

		if len(toDelete) != 0 {
			batch.Queue(stmtDeleteOldTagIDs, id, banner.FeatureID, toDelete)
		}

		if len(toInsert) != 0 {
			batch.Queue(stmtInsertNewTagIDs, id, banner.FeatureID, toInsert)
		}
	}

	if bannerPartial.FeatureID != nil {
		batch.Queue(stmtUpdateFeatureID, id, updatedBanner.FeatureID)
	}

	if bannerPartial.IsActive != nil || bannerPartial.Content != nil {
		newContentJSON, err := json.Marshal(updatedBanner.Content)
		if err != nil {
			return err
		}

		batch.Queue(stmtUpdateBanner, id, updatedBanner.IsActive, newContentJSON)
	}

	br := tx.SendBatch(ctx, batch)

	for i := 0; i != batch.Len(); i++ {
		ct, err := br.Exec()

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == SQLDuplicateErrCode {
				return service.ErrBannerAlreadyExists
			}
		}

		if ct.RowsAffected() == 0 {
			return fmt.Errorf(
				"err update banner with id=%d, rows_affected=%v",
				id,
				ct.RowsAffected(),
			)
		}
	}

	tx.Commit(ctx)
	return nil
}

func (repo BannerRepo) GetFiltered(ctx context.Context, filter bannermodels.FilterSchema) ([]bannermodels.Banner, error) {
	if !(filter.HasFeatureID || filter.HasTagID) {
		return repo.getFiltered(ctx, filter)
	}
	return repo.getFilteredWithFeatureAndTagFilter(ctx, filter)
}

func (repo BannerRepo) getFiltered(ctx context.Context, filter bannermodels.FilterSchema) ([]bannermodels.Banner, error) {
	var dbBanners []bannermodels.BannerDB
	err := pgxscan.Select(ctx, repo.db, &dbBanners, stmtBannerList, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	return bannermodels.SliceBannerDBToBanners(dbBanners)
}

func (repo BannerRepo) getFilteredWithFeatureAndTagFilter(ctx context.Context, filter bannermodels.FilterSchema) ([]bannermodels.Banner, error) {
	var stmtWhereFilter string
	args := []interface{}{
		filter.Limit,
		filter.Offset,
	}

	switch {
	case filter.HasFeatureID && filter.HasTagID:
		stmtWhereFilter = "feature_id=$3 AND tag_id=$4"
		args = append(args, filter.FeatureID, filter.TagID)
	case filter.HasFeatureID:
		stmtWhereFilter = "feature_id=$3"
		args = append(args, filter.FeatureID)
	case filter.HasTagID:
		args = append(args, filter.TagID)
		stmtWhereFilter = "tag_id=$3"
	default:
		return nil, fmt.Errorf(
			"at least one of the HasFeatureID and HasTagID should be true, got: %v, %v",
			filter.HasFeatureID,
			filter.HasTagID,
		)
	}

	stmtBannerListWithFilter := fmt.Sprintf(stmtBannerListWithFilterTemplate, stmtWhereFilter)

	var dbBanners []bannermodels.BannerDB
	err := pgxscan.Select(
		ctx,
		repo.db,
		&dbBanners,
		stmtBannerListWithFilter,
		args...,
	)
	if err != nil {
		return nil, err
	}

	return bannermodels.SliceBannerDBToBanners(dbBanners)
}
