package queryutils

import (
	"database/sql"
	"demogo/config"
	"demogo/pkg/logger"

	sq "github.com/elgris/sqrl"
)

var debug bool

func ToSql(q sq.Sqlizer, segment string) (string, []interface{}, error) {
	query, args, err := q.ToSql()
	if err != nil {
		return "", []interface{}{}, err
	}

	if debug {
		logger.Log.Infof(`%v QUERY: %v || ARGS: %v`, segment, query, args)
	}

	if err != nil {
		return "", []interface{}{}, err
	}

	return query, args, nil
}

func ErrRowsAffected(res sql.Result) error {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func New(conf *config.Config) {
	debug = conf.PostgreSQLDebug
}
