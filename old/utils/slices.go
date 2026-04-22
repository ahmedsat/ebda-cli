package utils

func FindF[T any](slice []T, predicate func(T) bool) T {
	for _, item := range slice {
		if predicate(item) {
			return item
		}
	}
	var zero T
	return zero
}

func UpdateF[T any](slice []T, predicate func(T) bool, updater func(T) T) []T {
	for i, item := range slice {
		if predicate(item) {
			slice[i] = updater(item)
		}
	}
	return slice
}

func TruePredicate[T any]() func(T) bool {
	return func(T) bool { return true }
}

func FindIndexF[T any](slice []T, predicate func(T) bool) int {
	for i, item := range slice {
		if predicate(item) {
			return i
		}
	}
	return -1
}
