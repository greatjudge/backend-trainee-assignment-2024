package models

type FilterSchema struct {
	HasFeatureID bool
	FeatureID    int

	HasTagID bool
	TagID    int

	Limit  int
	Offset int
}

func NewFilerSchema(limit int, offset int) FilterSchema {
	return FilterSchema{
		Limit:  limit,
		Offset: offset,
	}
}

func (fs *FilterSchema) SetFeatureID(featureID int) {
	fs.HasFeatureID = true
	fs.FeatureID = featureID
}

func (fs *FilterSchema) SetTagID(tagID int) {
	fs.HasTagID = true
	fs.TagID = tagID
}
