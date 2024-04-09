package tests

import (
	bannermodels "banner/internal/models/banner"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBannerList(t *testing.T) {
	db.SetUp(t, bannerTableName, bannerRelationTableName)
	defer db.TearDown(bannerTableName, bannerRelationTableName)

	// arrange
	banners := []bannermodels.Banner{
		{
			FeatureID: 1,
			TagIDs:    []int{1, 2},
			IsActive:  true,
			Content:   testContentObj,
		},
		{
			FeatureID: 1,
			TagIDs:    []int{3, 4},
			IsActive:  true,
			Content:   testContentObj,
		},
		{
			FeatureID: 2,
			TagIDs:    []int{1, 3},
			IsActive:  true,
			Content:   testContentObj,
		},
		{
			FeatureID: 3,
			TagIDs:    []int{1, 2, 3, 4, 5},
			IsActive:  true,
			Content:   testContentObj,
		},
		{
			FeatureID: 4,
			TagIDs:    []int{1, 2, 5, 7},
			IsActive:  true,
			Content:   testContentObj,
		},
	}

	bannersCreated, err := createBunners(banners)
	if err != nil {
		log.Panic(err)
	}

	url := bannerListURL + "?feature_id=1"

	client, req, err := makeClientRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panic(err)
	}

	// act
	resp, err := client.Do(req)

	// assert
	require.NoError(t, err, err)

	resultBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode, string(resultBytes))

	var resultBanners []bannermodels.Banner
	err = json.Unmarshal(resultBytes, &resultBanners)
	require.NoError(t, err, err)

	assert.Equal(t, 2, len(resultBanners), resultBanners)
	compareBanners(t, bannersCreated[0], resultBanners[1])
	compareBanners(t, bannersCreated[1], resultBanners[0])
}

func compareBanners(t *testing.T, expected, actual bannermodels.Banner) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.FeatureID, actual.FeatureID)
	assert.Equal(t, expected.TagIDs, actual.TagIDs)
	assert.Equal(t, expected.Content, actual.Content)
	assert.Equal(t, expected.IsActive, actual.IsActive)
}

// TODO add bad args, different fields
