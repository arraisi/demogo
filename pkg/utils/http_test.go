package utils

import (
	"github.com/arraisi/demogo/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetBaseURL(t *testing.T) {
	var (
		confHTTP = config.HTTPItemConfig{
			Host: "localhost",
			TLS:  false,
		}
		confHTTPS = config.HTTPItemConfig{
			Host: "localhost",
			TLS:  true,
		}
	)

	t.Run("success HTTP", func(t *testing.T) {
		baseURL := GetBaseURL(confHTTP)

		expectedHTTP := "http://localhost"

		require.Equal(t, expectedHTTP, baseURL)
	})

	t.Run("success HTTPS", func(t *testing.T) {
		baseURL := GetBaseURL(confHTTPS)

		expectedHTTPS := "https://localhost"

		require.Equal(t, expectedHTTPS, baseURL)
	})
}
