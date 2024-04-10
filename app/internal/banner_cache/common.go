package cache

import "fmt"

func formKeyFromTagIDFeatureID(tagID int, featureID int) string {
	return fmt.Sprintf("%d,%d", tagID, featureID)
}
