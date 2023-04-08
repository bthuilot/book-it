package util

func ListContains[T comparable](l []T, i T) (contains bool) {
	for _, li := range l {
		if li == i {
			contains = true
		}
	}
	return
}
