package utils

import (
	"database/sql"
	"fmt"
	"github.com/arraisi/demogo/config"
	"math"
	"reflect"
	"sort"
	"strings"
	"time"
)

// GeneratePostgreURL generates postgresql URL from DB Account struct
func GeneratePostgreURL(dbAcc config.DBAccount) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s  sslmode=disable extra_float_digits=-1 connect_timeout=%s", dbAcc.Username, dbAcc.Password, dbAcc.Name, dbAcc.Host, dbAcc.Port, dbAcc.Timeout)
	//return "user=" + dbAcc.Username + " password=" + dbAcc.Password + " dbname=" + dbAcc.DBName + " host=" + dbAcc.URL + " port=" + dbAcc.Port + "  sslmode=disable extra_float_digits=-1 connect_timeout=" + dbAcc.Timeout
}

func CalculateOffsetPagination(pageIndex int32, pageSize int32) int32 {
	pageIndex = pageIndex - 1
	if pageIndex < 1 {
		pageIndex = 0
	}
	return pageSize * pageIndex
}

func CalculateTotalPage(totalData int, pageSize int32) int {
	return int(math.Ceil(float64(totalData) / float64(pageSize)))
}

func GenerateMongoDBURL(account *config.MongoDBAccount) string {
	mongoHostArr := strings.Split(account.Host, ",")
	var allMongoURLHost []string
	for _, host := range mongoHostArr {
		allMongoURLHost = append(allMongoURLHost, fmt.Sprintf("%s:%s", host, account.Port))
	}
	allMongoURLHostStr := strings.Join(allMongoURLHost, ",")
	//local
	//return fmt.Sprintf("mongodb://%s:%s@%s/%s?ssl=false&authSource=admin", account.Username, account.Password, allMongoURLHostStr, account.DBName)
	return fmt.Sprintf("mongodb://%s:%s@%s/%s?ssl=false", account.Username, account.Password, allMongoURLHostStr, account.DBName)
}

func ArgQueryLikeContain(s string) string {
	return "%" + s + "%"
}

func ArgQueryLikeBegin(s string) string {
	return "%" + s
}

func ArgQueryLikeEnd(s string) string {
	return s + "%"
}

func GetTableColumns(i interface{}) (columns []string) {
	f := reflect.ValueOf(i)

	for i := 0; i < f.NumField(); i++ {
		tag := f.Type().Field(i).Tag.Get("db")
		if tag == "" {
			continue
		}

		columns = append(columns, tag)
	}

	return columns
}

func GetTableColumnWithAlias(i interface{}, alias string) (columns []string) {
	f := reflect.ValueOf(i)

	for i := 0; i < f.NumField(); i++ {
		tag := f.Type().Field(i).Tag.Get("db")
		if tag == "" {
			continue
		}

		columns = append(columns, fmt.Sprintf("%s.%s", alias, tag))
	}

	return columns
}

func GenerateColumnsAndValues(data map[string]interface{}) (columns []string, values []interface{}) {
	keys := make([]string, 0, len(data))

	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		columns = append(columns, k)
		values = append(values, data[k])
	}

	return columns, values
}

func ToSqlNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}

	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func ToSqlNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}
