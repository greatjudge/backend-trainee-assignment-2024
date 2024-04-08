package tools

// return slice with elems from lhs that is not in rhs
func SliceDiff(lhs []int, rhs []int) []int {
	result := make([]int, 0, len(lhs))
	rhsSet := SetFromSlice(rhs)

	for _, v := range lhs {
		if !rhsSet[v] {
			result = append(result, v)
		}
	}

	return result
}

func SetFromSlice(arr []int) map[int]bool {
	result := make(map[int]bool)
	for _, v := range arr {
		result[v] = true
	}
	return result
}
