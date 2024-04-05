package handler

import (
	"banner/internal/models"
	"banner/internal/sending"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type bannerService interface {
	GetUserBanner(ctx context.Context, tagID int, featureID int, useLastRevision bool, userToken string) (string, error)
	BannerList(ctx context.Context, filter FilterSchema) ([]models.Banner, error)
	CreateBanner(ctx context.Context, banner models.Banner) (int, error)
	DeleteBanner(ctx context.Context, id int) error
}

type BannerHandler struct {
	service bannerService
}

func NewBannerHandler(service bannerService) BannerHandler {
	return BannerHandler{
		service: service,
	}
}

func (h BannerHandler) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	tagID, err := tagIDFromQuery(queryParams)
	if err != nil {
		sending.SendErrorMsg(w, http.StatusBadRequest, badTagIDMsg)
		return
	}

	featureID, err := featureIDFromQuery(queryParams)
	if err != nil {
		sending.SendErrorMsg(w, http.StatusBadRequest, badFeatureIDMsg)
		return
	}

	useLastRevision := defaultUseLastVersion
	if queryParams.Has(useLastRevisionParamName) {
		useLastRevision, err = useLastRevisionFromQuery(queryParams)
		if err != nil {
			sending.SendErrorMsg(w, http.StatusBadRequest, badUseLastRevision)
		}
	}

	userToken, ok := r.Context().Value(UserTokenKey).(string)
	if !ok {
		sending.SendErrorMsg(w, http.StatusInternalServerError, errMsgUserTokenNotFound)
		return
	}

	bannerJSON, err := h.service.GetUserBanner(r.Context(), tagID, featureID, useLastRevision, userToken)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	sending.SendJSONBytes(w, http.StatusOK, []byte(bannerJSON))
}

func (h BannerHandler) BannerList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	limit := defaultLimit
	if queryParams.Has(limitParamName) {
		limitUint, err := StrToUint(queryParams.Get(limitParamName))
		if err != nil {
			sending.SendErrorMsg(w, http.StatusBadRequest, badLimitMsg)
			return
		}
		limit = int(limitUint)
	}

	offset := defaultOffset
	if queryParams.Has(offsetParamName) {
		offsetUint, err := StrToUint(queryParams.Get(offsetParamName))
		if err != nil {
			sending.SendErrorMsg(w, http.StatusBadRequest, badOfssetMsg)
			return
		}
		offset = int(offsetUint)
	}

	filter := NewFilerSchema(limit, offset)

	if queryParams.Has(featureIDParamName) {
		featureID, err := featureIDFromQuery(queryParams)
		if err != nil {
			sending.SendErrorMsg(w, http.StatusBadRequest, badFeatureIDMsg)
			return
		}
		filter.SetFeatureID(featureID)
	}

	if queryParams.Has(tagIDParamName) {
		tagID, err := tagIDFromQuery(queryParams)
		if err != nil {
			sending.SendErrorMsg(w, http.StatusBadRequest, badTagIDMsg)
			return
		}
		filter.SetTagID(tagID)
	}

	banners, err := h.service.BannerList(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	sending.JSONMarshallAndSend(w, http.StatusOK, banners)
}

func (h BannerHandler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sending.SendErrorMsg(w, http.StatusInternalServerError, errMsgCantReadBody)
		return
	}

	var bannerReq models.BannerRequest
	err = json.Unmarshal(body, &bannerReq)
	if err != nil {
		sending.SendErrorMsg(w, http.StatusBadRequest, err.Error()) // TODO
		return
	}

	id, err := h.service.CreateBanner(r.Context(), bannerReq.ToBanner())
	if err != nil {
		h.handleServiceError(w, err)
	}

	sending.JSONMarshallAndSend(w, http.StatusCreated, BannerIdMsg{ID: id})
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

func (h BannerHandler) UpdatePatial(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := IDFromVars(vars)
	if err != nil {
		sending.SendErrorMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		sending.SendErrorMsg(w, http.StatusInternalServerError, errMsgCantReadBody)
		return
	}

	// TODO
}

func (h BannerHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := IDFromVars(vars)
	if err != nil {
		sending.SendErrorMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.DeleteBanner(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h BannerHandler) handleServiceError(w http.ResponseWriter, err error) {
	// TODO
}
