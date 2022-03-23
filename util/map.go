package util

func MapWithErr[T1, T2 any](fromSlice []T1, mapFn func(from T1) (T2, error)) ([]T2, error) {
	var toSlice []T2

	for _, item := range fromSlice {
		mapped, err := mapFn(item)

		if err != nil {
			return toSlice, err
		}

		toSlice = append(toSlice, mapped)
	}

	return toSlice, nil
}
