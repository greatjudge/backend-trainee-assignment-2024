package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type BannerCreateRequest struct {
	TagIDs    []int                  `json:"tag_ids"`
	FeatureID int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
}

type BannerCreateResponse struct {
	BannerID int `json:"banner_id"`
}

func TestCreate(t *testing.T) {
	//arrange
	db.SetUp(t, bannerTableName, bannerRelationTableName)
	defer db.TearDown(bannerTableName, bannerRelationTableName)

	expectedID := 1

	bannerReq := BannerCreateRequest{
		TagIDs:    []int{1, 2, 3, 4, 5, 6},
		FeatureID: 1,
		Content:   testContentObj,
		IsActive:  true,
	}

	body, err := json.Marshal(bannerReq)
	if err != nil {
		log.Panic(err)
	}

	client, req, err := makeClientRequest(http.MethodPost, bannerCreateURL, bytes.NewBuffer(body))
	if err != nil {
		log.Panic(err)
	}

	// act
	resp, err := client.Do(req)

	// assert
	require.NoError(t, err, err)

	resultBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err, err)

	require.Equal(t, http.StatusCreated, resp.StatusCode, string(resultBytes))

	var b BannerCreateResponse
	err = json.Unmarshal(resultBytes, &b)
	require.NoError(t, err, err)

	assert.Equal(t, expectedID, b.BannerID)

	bannerInDB, err := getBannerByID(expectedID)
	require.NoError(t, err, err)

	assert.Equal(t, bannerReq.FeatureID, bannerInDB.FeatureID)
	assert.Equal(t, bannerReq.TagIDs, bannerInDB.TagIDs)
	assert.Equal(t, bannerReq.IsActive, bannerInDB.IsActive)
	assert.Equal(t, bannerReq.Content, bannerInDB.Content)
}

func TestCreateBadArgs(t *testing.T) {
	// TODO
}
