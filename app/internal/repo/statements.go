package repo

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

	stmtDeleteBanner = `
	DELETE from banner WHERE "id" = $1;
	DELETE from banner_relation WHERE banner_id = $1;
	`
)
