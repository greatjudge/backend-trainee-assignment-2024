package banner

import "errors"

var (
	ErrBadFeatureID = errors.New("bad type for featureID, expected int")
	ErrBadTagIDs    = errors.New("bad type for tagIDs, expected []int")
	ErrBadIsActive  = errors.New("bad type for isActive, expected bool")
	ErrBadContent   = errors.New("bad type for content, expected map[string]interface{}")
)
