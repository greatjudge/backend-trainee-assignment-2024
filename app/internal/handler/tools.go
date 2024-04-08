package handler

import (
	"errors"
	"net/url"
	"strconv"
)

const (
	tagIDParamName           = "tag_id"
	featureIDParamName       = "feature_id"
	useLastRevisionParamName = "use_last_revision"
	limitParamName           = "limit"
	offsetParamName          = "offset"
	idParamName              = "id"

	badTagIDMsg        = "tag_id должен быть целым числом"
	badTagIDsMsg       = "tag_ids должен быть массивом целых чисел"
	badContentMsg      = "content должен быть структурой"
	badFeatureIDMsg    = "feature_id должен быть целым числом"
	badIsActive        = "is_active должен быть типа bool"
	badUseLastRevision = "use_last_revision должен быть типа boolean"
	badLimitMsg        = "limit должен быть целым числом >= 0"
	badOfssetMsg       = "offset должен быть целым числом >= 0"
	badIDMsg           = "id должен быть целым числом"

	noIDinParamsMsg = "нужно указать id"

	errMsgUserNotFoundInCTX = "user not found in context"
	errMsgCantReadBody      = "can not read body"

	errMsgBannerNotFound      = "баннер не найден"
	errMsgBannerAlreadyExists = "баннер с такими feature_id и tag_id уже существует"

	defaultLimit          = 10
	defaultOffset         = 0
	defaultUseLastVersion = false
)

type BannerIdMsg struct {
	ID int `json:"banner_id"`
}

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

func IDFromVars(vars map[string]string) (int, error) {
	idStr, ok := vars[idParamName]
	if !ok {
		return 0, errors.New(noIDinParamsMsg)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New(badIDMsg)
	}

	return id, nil
}
