package v1beta1

import "strings"

func EventKey(bits ...string) string {
	key := []string{"processor", "v1beta1"}
	key = append(key, bits...)
	return strings.Join(key, ".")
}

func CacheKey(bits ...string) string {
	key := []string{"processor", "v1beta1"}
	key = append(key, bits...)
	return strings.Join(key, ":")
}
