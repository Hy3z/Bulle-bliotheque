package util

func CastArray[T any](array []any) []T {
	var tarray []T
	for _, v := range array {
		val, ok := v.(T)
		if ok {
			tarray = append(tarray, val)
		}
	}
	return tarray
}
