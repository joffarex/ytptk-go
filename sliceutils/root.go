package sliceutils

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func Filter[T any](ts []T, f func(T) bool) []T {
	var arr []T
	for _, item := range ts {
		if f(item) {
			arr = append(arr, item)
		}
	}
	return arr
}
