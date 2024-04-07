package service

import "errors"

var (
	ErrUserForbidden = errors.New("the user does not have access")

	ErrBannerNotFound      = errors.New("banner not found")
	ErrBannerAlreadyExists = errors.New(
		"banner with this tag_ids and feature_id already exists",
	)

	ErrDBBannerNotFound      = errors.New("banner not found in db")
	ErrDBBannerAlreadyExists = errors.New(
		"banner with this tag_ids and feature_id already exists",
	)

	ErrCacheBannerNotFound = errors.New("banner not found in cache")
)
