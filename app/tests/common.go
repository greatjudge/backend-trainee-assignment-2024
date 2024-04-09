package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	bannermodels "banner/internal/models/banner"
)

const (
	baseURL          = "http://localhost:9000"
	bannerCreateURL  = baseURL + "/banner"
	bannerListURL    = baseURL + "/banner"
	bannerGetUserURL = baseURL + "/user_banner"
	bannerUpdateURL  = baseURL + "/banner/%d"
	bannerDeleteURL  = baseURL + "/banner/%d"

	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"

	bannerTableName         = "banner"
	bannerRelationTableName = "banner_relation"

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

	stmtCreateBanner = `
	with create_banner AS (
		INSERT into banner (feature_id, is_active, "content") VALUES ($2, $3, $4) RETURNING "id", feature_id
	),
	create_banner_relation as (
		INSERT into banner_relation (banner_id, feature_id, tag_id)
		SELECT cb.id as banner_id, cb.feature_id as feature_id, UNNEST($1::int[]) as tag_id FROM create_banner AS cb
	)
	  
	SELECT "id" FROM create_banner;
	`
)

var (
	// testContent    = []byte(`{"title": "some_title", "text": "some_text", "url": "some_url"}`)
	testContentObj = map[string]interface{}{"title": "some_title", "text": "some_text", "url": "some_url"}
)

func makeClient() *http.Client {
	return &http.Client{}
}

func makeClientRequest(
	method string,
	url string,
	body io.Reader,
) (*http.Client, *http.Request, error) {
	client := makeClient()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add(contentTypeHeader, contentTypeJSON)

	return client, req, nil
}

func getBannerByID(id int) (bannermodels.Banner, error) {
	var contentJSON []byte
	var banner bannermodels.Banner

	row := db.DB.QueryRow(context.Background(), stmtGetBannerByID, id)

	err := row.Scan(
		&banner.ID,
		&banner.FeatureID,
		&banner.TagIDs,
		&contentJSON,
		&banner.IsActive,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)

	if err != nil {
		return bannermodels.Banner{}, err
	}

	err = json.Unmarshal(contentJSON, &banner.Content)
	if err != nil {
		return bannermodels.Banner{}, err
	}

	return banner, nil
}

func createBanner(banner bannermodels.Banner) (bannermodels.Banner, error) {
	contentJSON, err := json.Marshal(banner.Content)
	if err != nil {
		return bannermodels.Banner{}, err
	}

	row := db.DB.QueryRow(
		context.Background(),
		stmtCreateBanner,
		banner.TagIDs,
		banner.FeatureID,
		banner.IsActive,
		contentJSON,
	)

	var id int
	err = row.Scan(&id)

	if err != nil {
		return bannermodels.Banner{}, err
	}

	banner.ID = id

	return banner, nil
}

func createBunners(banners []bannermodels.Banner) ([]bannermodels.Banner, error) {
	result := make([]bannermodels.Banner, len(banners))
	for i, b := range banners {
		b, err := createBanner(b)
		if err != nil {
			return nil, err
		}
		result[i] = b
	}
	return result, nil
}
