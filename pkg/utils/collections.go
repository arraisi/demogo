package utils

import (
	"fmt"
	"strconv"
)

// GetUniqueElementsSliceOfString gets unique element of a slice of string
func GetUniqueElementsSliceOfString(collection []string) []string {
	result := make([]string, 0)
	hashMap := make(map[string]int)

	for _, element := range collection {
		if _, oke := hashMap[element]; !oke {
			hashMap[element] = 1
		}
	}

	for key := range hashMap {
		result = append(result, key)
	}
	return result
}

// ConvertSliceOfStringIntoSliceOfInt64 converts slice of string into slice of int64
func ConvertSliceOfStringIntoSliceOfInt64(collection []string) []int64 {
	result := make([]int64, len(collection))

	for id, element := range collection {
		result[id], _ = strconv.ParseInt(element, 10, 64)
	}
	return result
}

// ConstructSliceOfInt64IntoString create a string from slice of int64 with a custom separator
func ConstructSliceOfInt64IntoString(collection []int64, separator string) string {
	result := ""
	for i := 0; i < len(collection); i++ {
		if i == 0 {
			result = fmt.Sprintf("%d", collection[i])
		} else {
			result = fmt.Sprintf("%s%s%d", result, separator, collection[i])
		}
	}
	return result
}
