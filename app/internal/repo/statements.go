package repo

// TODO add update_at set

const (
	SQLDuplicateErrCode = "23505"

	stmtCreateBanner = `
	with create_banner AS (
		INSERT into banner (feature_id, is_active, "content") VALUES ($2, $3, $4) RETURNING "id", feature_id
	),
	create_banner_relation as (
		INSERT into banner_relation (banner_id, feature_id, tag_id)
		SELECT cb.id as banner_id
		       , cb.feature_id as feature_id
		       , UNNEST($1::int[]) as tag_id 
		  FROM create_banner AS cb
	)
	  
	SELECT "id" FROM create_banner;
	`

	stmtGetUserBanner = `
	with find_banner as (
		SELECT banner_id, tag_id, feature_id FROM banner_relation WHERE feature_id=$2 AND tag_id=$1
	)
	
	SELECT
		b.id,
		b.feature_id,
		(SELECT ARRAY_AGG(tag_id) FROM banner_relation AS br WHERE br.banner_id = b.id) as tag_ids,
		b.content,
		b.is_active,
		b.created_at,
		b.updated_at
	FROM banner as b JOIN find_banner as fb ON (b.id = fb.banner_id);
	`

	stmtGetBannerByID = `
	SELECT
		b.id,
		b.feature_id,
		(SELECT ARRAY_AGG(tag_id) FROM banner_relation AS br WHERE br.banner_id = b.id) as tag_ids,
		b.content,
		b.is_active,
		b.created_at,
		b.updated_at
	FROM banner as b
	WHERE b.id = $1;
	`

	stmtUpdateBannerTemplate = `
	UPDATE banner SET %v WHERE "id"=$1;
	`

	stmtDeleteOldTagIDs = `
	DELETE from banner_relation WHERE banner_id=$1 AND feature_id=$2 AND tag_id = ANY($3::int[]);
	`

	stmtInsertNewTagIDs = `
	INSERT INTO banner_relation (banner_id, feature_id, tag_id) SELECT $1, $2, UNNEST($3::int[]);
	`

	stmtUpdateFeatureID = `
	UPDATE banner_relation SET feature_id=$2 WHERE banner_id=$1;
	`

	stmtBannerList = `
	SELECT
		b.id,
		b.feature_id,
		(SELECT ARRAY_AGG(tag_id) FROM banner_relation AS br WHERE br.banner_id = b.id) as tag_ids,
		b.content,
		b.is_active,
		b.created_at,
		b.updated_at
	FROM banner as b
	ORDER BY b.created_at DESC
	LIMIT $1 OFFSET $2
	`

	stmtBannerListWithFilterTemplate = `
	with filtered_banners as (
	  SELECT DISTINCT banner_id FROM banner_relation WHERE %v
	)
	
	SELECT
		b.id,
		b.feature_id,
		(SELECT ARRAY_AGG(tag_id) FROM banner_relation AS br WHERE br.banner_id = b.id) as tag_ids,
		b.content,
		b.is_active,
		b.created_at,
		b.updated_at
	FROM banner as b JOIN filtered_banners as fb ON (b.id = fb.banner_id)
	ORDER BY b.created_at DESC
	LIMIT $1 OFFSET $2;
	`

	stmtDeleteBanner = `
	DELETE from banner WHERE "id" = $1;
	`
)
