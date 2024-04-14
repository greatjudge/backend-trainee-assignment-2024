package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	bannermodels "banner/internal/models/banner"
	"banner/internal/service"
	"banner/internal/tools"
)

type database interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type BannerRepo struct {
	db database
}

func NewBannerRepo(db database) *BannerRepo {
	return &BannerRepo{
		db: db,
	}
}

func (repo *BannerRepo) CreateBanner(ctx context.Context, banner bannermodels.Banner) (int, error) {
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
		banner.IsActive,
		contentJSON,
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

func (repo *BannerRepo) GetUserBanner(ctx context.Context, tagID int, featureID int) (bannermodels.Banner, error) {
	// ADD TRANSATION ?
	row := repo.db.QueryRow(ctx, stmtGetUserBanner, tagID, featureID)

	var contentJSON []byte
	var banner bannermodels.Banner

	err := row.Scan(
		&banner.ID,
		&banner.FeatureID,
		&banner.TagIDs,
		&contentJSON,
		&banner.IsActive,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return bannermodels.Banner{}, service.ErrDBBannerNotFound
	case err != nil:
		return bannermodels.Banner{}, err
	}

	err = json.Unmarshal(contentJSON, &banner.Content)
	if err != nil {
		return bannermodels.Banner{}, err
	}

	return banner, nil
}

func (repo *BannerRepo) PartialUpdateBanner(ctx context.Context, id int, bannerPartial bannermodels.BannerPartialUpdate) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, stmtGetBannerByID, id)

	var contentJSON []byte
	var banner bannermodels.Banner

	err = row.Scan(
		&banner.ID,
		&banner.FeatureID,
		&banner.TagIDs,
		&contentJSON,
		&banner.IsActive,
		&banner.CreatedAt,
		&banner.UpdatedAt,
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

	updateArgs := []interface{}{id}
	updateArgsExtend, updateFields, err := repo.formUpdateArgsFields(
		bannerPartial,
		updatedBanner,
		len(updateArgs)+1,
	)
	if err != nil {
		return err
	}

	if len(updateArgsExtend) != 0 {
		updateArgs = append(updateArgs, updateArgsExtend...)

		updateSetString := strings.Join(updateFields, ", ")
		stmtUpdateBanner := fmt.Sprintf(stmtUpdateBannerTemplate, updateSetString)

		batch.Queue(stmtUpdateBanner, updateArgs...)
	}

	br := tx.SendBatch(ctx, batch)

	for i := 0; i != batch.Len(); i++ {
		ct, err := br.Exec()

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == SQLDuplicateErrCode {
				return service.ErrBannerAlreadyExists
			}
			return err
		}

		if ct.RowsAffected() == 0 {
			return fmt.Errorf(
				"err update banner with id=%d, rows_affected=%v must be >= 1",
				id,
				ct.RowsAffected(),
			)
		}
	}

	err = br.Close()
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BannerRepo) GetFiltered(ctx context.Context, filter bannermodels.FilterSchema) ([]bannermodels.Banner, error) {
	if !(filter.HasFeatureID || filter.HasTagID) {
		return repo.getFiltered(ctx, filter)
	}
	return repo.getFilteredWithFeatureAndTagFilter(ctx, filter)
}

func (repo *BannerRepo) getFiltered(ctx context.Context, filter bannermodels.FilterSchema) ([]bannermodels.Banner, error) {
	var dbBanners []bannermodels.BannerDB
	err := repo.db.Select(ctx, &dbBanners, stmtBannerList, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	return bannermodels.SliceBannerDBToBanners(dbBanners)
}

func (repo *BannerRepo) getFilteredWithFeatureAndTagFilter(ctx context.Context, filter bannermodels.FilterSchema) ([]bannermodels.Banner, error) {
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
	err := repo.db.Select(
		ctx,
		&dbBanners,
		stmtBannerListWithFilter,
		args...,
	)
	if err != nil {
		return nil, err
	}

	return bannermodels.SliceBannerDBToBanners(dbBanners)
}

func (repo *BannerRepo) DeleteBanner(ctx context.Context, id int) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	ct, err := tx.Exec(ctx, stmtDeleteBanner, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return service.ErrDBBannerNotFound
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *BannerRepo) formUpdateArgsFields(
	bannerPartial bannermodels.BannerPartialUpdate,
	updatedBanner bannermodels.Banner,
	nextArgnum int,
) ([]interface{}, []string, error) {

	updateArgs := make([]interface{}, 0, 3)
	updateFields := make([]string, 0, 3)

	if bannerPartial.TagIDs != nil {
		updateArgs = append(updateArgs, updatedBanner.TagIDs)
		updateFields = append(
			updateFields,
			fmt.Sprintf("tag_ids = $%d", nextArgnum), // TODO move str to const
		)
		nextArgnum += 1
	}

	if bannerPartial.FeatureID != nil {
		updateArgs = append(updateArgs, updatedBanner.FeatureID)
		updateFields = append(
			updateFields,
			fmt.Sprintf("feature_id = $%d", nextArgnum), // TODO move str to const
		)
		nextArgnum += 1
	}

	if bannerPartial.IsActive != nil {
		updateArgs = append(updateArgs, bannerPartial.IsActive)
		updateFields = append(
			updateFields,
			fmt.Sprintf("is_active = $%d", nextArgnum), // TODO move str to const
		)
		nextArgnum += 1
	}

	if bannerPartial.Content != nil {
		newContentJSON, err := json.Marshal(updatedBanner.Content)
		if err != nil {
			return nil, nil, err
		}

		updateArgs = append(updateArgs, newContentJSON)
		updateFields = append(
			updateFields,
			fmt.Sprintf(`"content" = $%d`, nextArgnum), // TODO move str to const
		)
		nextArgnum += 1
	}

	return updateArgs, updateFields, nil
}
