package service

import "errors"

var (
	ErrBannerNotFound = errors.New("banner not found")
	ErrUserForbidden  = errors.New("the user does not have access")
)
