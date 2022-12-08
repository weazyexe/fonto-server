package slices

func Any[T interface{}](slice []T, predicate func(T) bool) bool {
	for _, elem := range slice {
		if predicate(elem) {
			return true
		}
	}
	return false
}
