package tests

import (
	bannermodels "banner/internal/models/banner"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteBanner(t *testing.T) {
	db.SetUp(t, bannerTableName, bannerRelationTableName)
	defer db.TearDown(bannerTableName, bannerRelationTableName)

	// arrange
	banner := bannermodels.Banner{
		FeatureID: 1,
		TagIDs:    []int{1, 2},
		IsActive:  true,
		Content:   testContentObj,
	}

	banner, err := createBanner(banner)
	if err != nil {
		log.Panic(err)
	}

	url := fmt.Sprintf(bannerDeleteURL, banner.ID)
	client, req, err := makeClientRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Panic(err)
	}

	// act
	resp, err := client.Do(req)

	// assert
	require.NoError(t, err, err)

	resultBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err, err)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode, string(resultBytes))

	_, err = getBannerByID(banner.ID)
	require.Error(t, err)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}
