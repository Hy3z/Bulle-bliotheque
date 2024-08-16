package util

// CastArray convertit directement un type inconnu en un tableau d'un type défini
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
