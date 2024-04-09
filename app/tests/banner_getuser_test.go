package tests

import (
	bannermodels "banner/internal/models/banner"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserBanner(t *testing.T) {
	db.SetUp(t, bannerTableName, bannerRelationTableName)
	defer db.TearDown(bannerTableName, bannerRelationTableName)

	// arrange
	contentObj := testContentObj

	tag_id, feature_id := 1, 1
	banner := bannermodels.Banner{
		TagIDs:    []int{tag_id},
		FeatureID: feature_id,
		Content:   contentObj,
		IsActive:  true,
	}

	_, err := createBanner(banner)
	if err != nil {
		log.Panic(err)
	}

	url := bannerGetUserURL + fmt.Sprintf("?tag_id=%v&feature_id=%v", tag_id, feature_id)

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

	var resultContentObj map[string]interface{}
	err = json.Unmarshal(resultBytes, &resultContentObj)

	require.NoError(t, err, string(resultBytes))
	assert.Equal(t, contentObj, resultContentObj)
}

// TODO BAD ARGS
