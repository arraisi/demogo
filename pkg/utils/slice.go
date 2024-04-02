package utils

// ToMap convert array to map
func ToMap[K comparable, V any](arr []V, f func(V) (K, V)) map[K]V {
	m := make(map[K]V)
	for _, v := range arr {
		key, val := f(v)
		m[key] = val
	}
	return m
}

// Mapper transform array of T to array of U
func Mapper[T, U any](arr []T, f func(T) U) []U {
	res := make([]U, len(arr))
	for i, v := range arr {
		res[i] = f(v)
	}
	return res
}
