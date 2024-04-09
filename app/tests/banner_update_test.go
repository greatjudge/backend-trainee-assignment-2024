package tests

import (
	bannermodels "banner/internal/models/banner"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateBanner(t *testing.T) {
	db.SetUp(t, bannerTableName, bannerRelationTableName)
	defer db.TearDown(bannerTableName, bannerRelationTableName)

	// arrange
	contentObj := testContentObj
	oldContentObj := map[string]interface{}{
		"title": "new title",
		"text":  "new text",
		"url":   "new url",
	}

	banner := bannermodels.Banner{
		TagIDs:    []int{1, 2, 3, 4, 5},
		FeatureID: 1,
		Content:   contentObj,
		IsActive:  true,
	}

	banner, err := createBanner(banner)
	if err != nil {
		log.Panic(err)
	}

	newBanner := bannermodels.Banner{
		ID:        banner.ID,
		TagIDs:    []int{1, 2, 3, 6},
		FeatureID: 2,
		Content:   oldContentObj,
		IsActive:  false,
	}

	bodyObj := map[string]interface{}{
		"tag_ids":    newBanner.TagIDs,
		"feature_id": newBanner.FeatureID,
		"content":    newBanner.Content,
		"is_active":  newBanner.IsActive,
	}

	body, err := json.Marshal(bodyObj)
	if err != nil {
		log.Panic(err)
	}

	url := fmt.Sprintf(bannerUpdateURL, banner.ID)

	client, req, err := makeClientRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		log.Panic(err)
	}

	// act
	resp, err := client.Do(req)

	// assert
	require.NoError(t, err, err)

	resultBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err, err)

	require.Equal(t, http.StatusOK, resp.StatusCode, string(resultBytes))

	bannerInDB, err := getBannerByID(banner.ID)
	require.NoError(t, err, err)

	assert.Equal(t, newBanner.FeatureID, bannerInDB.FeatureID)
	assert.Equal(t, newBanner.TagIDs, bannerInDB.TagIDs)
	assert.Equal(t, newBanner.IsActive, bannerInDB.IsActive)
	assert.Equal(t, newBanner.Content, bannerInDB.Content)
}

// TODO bad args; some args (not all); empty
