package utils

import (
	"context"
	"demogo/pkg/constant"
	"demogo/pkg/utils/jsoniter"
	"errors"
	"fmt"
	"math"
)

func RemoveSliceByIndex[T any](s []T, index int) []T {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

func AbsInt[T int | int64 | int32](num T) T {
	if num < 0 {
		return -num
	}

	return num
}

func GetFloat64Abs[T int | int64 | int32](input T) float64 {
	return math.Abs(float64(input))
}

// JSONString returns JSON string of data, if error occurs, it will return fmt.Sprintf("%#v", data).
// This function is used to print data in log.
// Don't use if you want to unmarshal the data later.
func JSONString(data interface{}) string {
	str, err := jsoniter.MarshalToString(data)
	if err != nil {
		return fmt.Sprintf("%#v", data)
	}

	return str
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MinimumOfArrInt(arr []int64) (int64, error) {
	if len(arr) == 0 {
		return 0, errors.New("empty slice")
	}

	min := arr[0]
	for _, i := range arr {
		if i < min {
			min = i
		}
	}

	return min, nil
}

func RemoveStringFromArray(arr *[]string, strToRemove string) {
	var idx int
	for _, str := range *arr {
		if str != strToRemove {
			(*arr)[idx] = str
			idx++
		}
	}
	*arr = (*arr)[:idx]
}

func GetRequestIDFromContext(ctx context.Context) (result string) {
	ctxValue := ctx.Value(constant.REQUEST_ID)
	if ctxValue == nil {
		return result
	}

	result = ctxValue.(string)
	return result
}
