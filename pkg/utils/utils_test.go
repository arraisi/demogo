package utils

import (
	"context"
	"errors"
	"github.com/arraisi/demogo/pkg/constant"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUtil_GetRequestIDFromContext(t *testing.T) {
	parentContext := context.Background()
	testCases := []struct {
		name           string
		ctx            context.Context
		expectedResult string
	}{
		{
			name:           "success case with context value",
			ctx:            context.WithValue(parentContext, constant.REQUEST_ID, "1"),
			expectedResult: "1",
		},
		{
			name:           "success case with empty value",
			ctx:            parentContext,
			expectedResult: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult := GetRequestIDFromContext(tc.ctx)
			if actualResult != tc.expectedResult {
				t.Errorf("Expected result and actual result are not equal!\n expected : %s\n actual : %s\n", tc.expectedResult, actualResult)
			}
		})
	}
}

func TestUtil_StringInSlice(t *testing.T) {
	t.Run("should return true", func(t *testing.T) {
		// given
		slice := []string{"a", "b", "c"}

		// when
		result := StringInSlice("a", slice)

		// then
		require.True(t, result)
	})

	t.Run("should return false", func(t *testing.T) {
		// given
		slice := []string{"a", "b", "c"}

		// when
		result := StringInSlice("d", slice)

		// then
		require.False(t, result)
	})
}

func TestUtil_MinimumOfArrInt(t *testing.T) {
	testCases := []struct {
		name           string
		requestParam   []int64
		expectedResult int64
		expectedError  error
	}{
		{
			name:           "success get minimum value",
			requestParam:   []int64{12, 6, 32, 54},
			expectedResult: 6,
			expectedError:  nil,
		},
		{
			name:           "empty slice",
			requestParam:   []int64{},
			expectedResult: 0,
			expectedError:  errors.New("empty slice"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult, actualErr := MinimumOfArrInt(tc.requestParam)
			require.Equal(t, actualResult, tc.expectedResult)
			require.Equal(t, actualErr, tc.expectedError)
		})
	}
}

func TestRemoveStringFromArray(t *testing.T) {
	arr := []string{"apple", "banana", "orange", "banana", "grape"}
	strToRemove := "banana"

	RemoveStringFromArray(&arr, strToRemove)

	expected := []string{"apple", "orange", "grape"}
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("Expected %v, but got %v", expected, arr)
	}
}
