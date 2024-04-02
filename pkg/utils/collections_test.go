package utils

import (
	"reflect"
	"testing"
)

var emptySlice = "empty slice"

func TestUtil_GetUniqueElementsSliceOfString(t *testing.T) {
	testCases := []struct {
		name           string
		inputSlice     []string
		expectedResult []string
	}{
		{
			name:           "normal slice duplicate",
			inputSlice:     []string{"this", "one", "and", "that", "one", "are", "yours", "and", "shut", "up"},
			expectedResult: []string{"this", "one", "and", "that", "are", "yours", "shut", "up"},
		},
		{
			name:           emptySlice,
			inputSlice:     []string{},
			expectedResult: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult := GetUniqueElementsSliceOfString(tc.inputSlice)
			if len(actualResult) != len(tc.expectedResult) {
				t.Error("slice length are not the same")
			}
		})
	}
}

func TestUtil_ConvertSliceOfStringIntoSliceOfInt64(t *testing.T) {
	testCases := []struct {
		name           string
		rawSlice       []string
		expectedResult []int64
	}{
		{
			name:           "normal slice",
			rawSlice:       []string{"1", "2", "3", "666"},
			expectedResult: []int64{1, 2, 3, 666},
		},
		{
			name:           emptySlice,
			rawSlice:       []string{},
			expectedResult: []int64{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult := ConvertSliceOfStringIntoSliceOfInt64(tc.rawSlice)
			if len(actualResult) != len(tc.expectedResult) {
				t.Error("slice length are not the same")
			}
			if !reflect.DeepEqual(actualResult, tc.expectedResult) {
				t.Error("actual and expected slice are not the same")
			}
		})
	}
}

func TestUtil_ConstructSliceOfInt64IntoString(t *testing.T) {
	testCases := []struct {
		name           string
		rawSlice       []int64
		separator      string
		expectedResult string
	}{
		{
			name:           "normal slice",
			rawSlice:       []int64{1, 2, 3, 4, 5, 6},
			separator:      ",",
			expectedResult: "1,2,3,4,5,6",
		},
		{
			name:           emptySlice,
			rawSlice:       []int64{},
			separator:      "-",
			expectedResult: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult := ConstructSliceOfInt64IntoString(tc.rawSlice, tc.separator)
			if actualResult != tc.expectedResult {
				t.Error("actual result and expected result are not the same")
			}
			if !reflect.DeepEqual(actualResult, tc.expectedResult) {
				t.Error("actual and expected slice are not the same")
			}
		})
	}
}
