package queryutils

import (
	"github.com/arraisi/demogo/config"
	"github.com/arraisi/demogo/pkg/logger"
	"log"
	"testing"

	sq "github.com/elgris/sqrl"
	"github.com/stretchr/testify/assert"
)

var (
	segmentTest = "test"
	mockConf    = &config.Config{
		PostgreSQLDebug: false,
	}
)

func TestQueryer(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		New(mockConf)
		assert.Equal(t, debug, false)
	})

	t.Run("ToSql select", func(t *testing.T) {
		New(mockConf)
		_, _, err := ToSql(sq.Select("*"), segmentTest)
		assert.Nil(t, err)
	})

	t.Run("ToSql update", func(t *testing.T) {
		New(mockConf)
		q := sq.Update("mock").Set("mock = ?", "mock")
		_, _, err := ToSql(q, segmentTest)
		assert.Nil(t, err)
	})

	t.Run("ToSql insert", func(t *testing.T) {
		New(mockConf)
		q := sq.Insert("mock").Values("mock = ?", "mock")
		_, _, err := ToSql(q, segmentTest)
		assert.Nil(t, err)
	})

	t.Run("ToSql delete", func(t *testing.T) {
		New(mockConf)
		q := sq.Delete("mock")
		_, _, err := ToSql(q, segmentTest)
		assert.Nil(t, err)
	})

	t.Run("ToSql error", func(t *testing.T) {
		_, _, err := ToSql(&sq.SelectBuilder{}, segmentTest)
		assert.Error(t, err)
	})

	t.Run("ToSql with debug", func(t *testing.T) {
		customLog, err := logger.InitLogger("info")
		if err != nil {
			log.Fatalf("init logging failed: %v", err)
		}
		logger.SetLogger(customLog)

		mockConfDebug := &config.Config{
			PostgreSQLDebug: true,
		}

		New(mockConfDebug)
		_, _, err = ToSql(sq.Select("*"), segmentTest)
		assert.Nil(t, err)
	})

}
