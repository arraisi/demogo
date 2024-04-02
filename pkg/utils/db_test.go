package utils

import (
	"testing"
)

func TestUtil_GeneratePostgreURL(t *testing.T) {
	testCases := []struct {
		name           string
		account        config.DBAccount
		expectedResult string
	}{
		{
			name: "success case",
			account: config.DBAccount{
				Username: "admin",
				Password: "admin",
				URL:      "localhost",
				Port:     "5432",
				DBName:   "inventory",
				Timeout:  "10",
			},
			expectedResult: "user=admin password=admin dbname=inventory host=localhost port=5432  sslmode=disable extra_float_digits=-1 connect_timeout=10",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult := GeneratePostgreURL(tc.account)
			if actualResult != tc.expectedResult {
				t.Errorf("Expected result and actual result are not equal!\n expected : %s\n actual : %s\n", tc.expectedResult, actualResult)
			}
		})
	}
}

func TestUtil_CalculateOffsetPagination(t *testing.T) {
	successCalculate := "success calculate"
	testCases := []struct {
		name           string
		pageIndex      int32
		pageSize       int32
		expectedResult int32
	}{
		{
			name:           successCalculate,
			pageIndex:      0,
			pageSize:       20,
			expectedResult: int32(0),
		},
		{
			name:           successCalculate,
			pageIndex:      1,
			pageSize:       20,
			expectedResult: int32(0),
		},
		{
			name:           successCalculate,
			pageIndex:      2,
			pageSize:       20,
			expectedResult: int32(20),
		},
		{
			name:           successCalculate,
			pageIndex:      10,
			pageSize:       20,
			expectedResult: int32(180),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult := CalculateOffsetPagination(tc.pageIndex, tc.pageSize)
			if actualResult != tc.expectedResult {
				t.Errorf("Expected result and actual result are not equal!\n expected : %d\n actual : %d\n", tc.expectedResult, actualResult)
			}
		})
	}
}

func TestUtil_CalculateTotalPage(t *testing.T) {
	testCases := []struct {
		name           string
		totalData      int
		pageSize       int32
		expectedResult int
	}{
		{
			name:           "divisible by the numerator",
			totalData:      100,
			pageSize:       20,
			expectedResult: 5,
		},
		{
			name:           "still have page after dividing",
			totalData:      85,
			pageSize:       20,
			expectedResult: 5,
		},
		{
			name:           "still have page after dividing",
			totalData:      85,
			pageSize:       10,
			expectedResult: 9,
		},
		{
			name:           "expect zero",
			totalData:      0,
			pageSize:       20,
			expectedResult: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult := CalculateTotalPage(tc.totalData, tc.pageSize)
			if actualResult != tc.expectedResult {
				t.Errorf("Expected result and actual result are not equal!\n expected : %d\n actual : %d\n", tc.expectedResult, actualResult)
			}
		})
	}
}
