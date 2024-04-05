package handler

import (
	"net/url"
	"strconv"
)

type userTokenT string

const UserTokenKey userTokenT = "user token"

const (
	tagIDParamName           = "tag_id"
	featureIDParamName       = "feature_id"
	useLastRevisionParamName = "use_last_revision"
	limitParamName           = "limit"
	offsetParamName          = "offset"
	idParamName              = "id"

	badTagIDMsg        = "tag_id должен быть целым числом"
	badFeatureIDMsg    = "feature_id должен быть целым числом"
	badUseLastRevision = "use_last_revision должен быть типа boolean"
	badLimitMsg        = "limit должен быть целым числом >= 0"
	badOfssetMsg       = "offset должен быть целым числом >= 0"
	badIDMsg           = "id должен быть целым числом"

	noIDinParamsMsg = "нужно указать id"

	errMsgUserTokenNotFound = "user_token not found in context"
	errMsgCantReadBody      = "can not read body"

	defaultLimit          = 10
	defaultOffset         = 0
	defaultUseLastVersion = false
)

func featureIDFromQuery(queryParams url.Values) (int, error) {
	return strconv.Atoi(queryParams.Get(tagIDParamName))
}

func tagIDFromQuery(queryParams url.Values) (int, error) {
	return strconv.Atoi(queryParams.Get(tagIDParamName))
}

func useLastRevisionFromQuery(queryParams url.Values) (bool, error) {
	return strconv.ParseBool(queryParams.Get(useLastRevisionParamName))
}

func StrToUint(str string) (uint, error) {
	val, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

type BannerIdMsg struct {
	ID int `json:"banner_id"`
}