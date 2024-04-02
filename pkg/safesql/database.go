package safesql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/arraisi/demogo/config"
	"github.com/arraisi/demogo/pkg/utils"
	"log"
	"time"

	"github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

const (
	defaultTimeOutString = "5"
)

// Database struct
type PostgreSQLHandler struct {
	DBRead  *sqlx.DB
	DBWrite *sqlx.DB
	Tx      *sqlx.Tx
}

// SqlxTx is a wrapper struct for sqlx Tx
type SqlxTx struct {
	Tx *sqlx.Tx
}

func (d *SqlxTx) Commit() {
	_ = d.Tx.Commit()
}

func (d *SqlxTx) Rollback() {
	_ = d.Tx.Rollback()
}

func (d *PostgreSQLHandler) ConnectDB(
	dbAccRead *config.DBAccount,
	dbAccWrite *config.DBAccount,
) {

	if dbAccRead.Timeout == "" {
		dbAccRead.Timeout = defaultTimeOutString
	}

	sqltrace.Register("postgres", &pq.Driver{}, sqltrace.WithServiceName("postgresql-ims"))
	dbRead, err := sqlxtrace.Open("postgres", utils.GeneratePostgreURL(*dbAccRead))
	if err != nil {
		log.Fatalf("Failed to open connection to DB Read! Error : %s\n", err.Error())
	}

	d.DBRead = dbRead
	err = d.DBRead.Ping()
	if err != nil {
		log.Fatalf("Failed to test connection (PING) to DB Read! Error : %s\n", err.Error())
	}

	log.Println("Successfully connect to DB Read")
	if dbAccWrite.Timeout == "" {
		dbAccWrite.Timeout = defaultTimeOutString
	}

	sqltrace.Register("postgres", &pq.Driver{}, sqltrace.WithServiceName("postgresql-ims"))
	dbWrite, err := sqlxtrace.Open("postgres", utils.GeneratePostgreURL(*dbAccWrite))
	if err != nil {
		log.Fatalf("Failed to open connection to DB Write! Error : %s\n", err.Error())
	}

	d.DBWrite = dbWrite
	err = d.DBWrite.Ping()
	if err != nil {
		log.Fatalf("Failed to test connection (PING) to Db Write! Error : %v\n", err.Error())
	}

	log.Println("Successfully connect to DB Write")

	d.DBWrite.SetConnMaxLifetime(time.Duration(dbAccWrite.MaxLifeTime) * time.Second)
	d.DBRead.SetConnMaxLifetime(time.Duration(dbAccRead.MaxLifeTime) * time.Second)

	// max connection
	d.DBWrite.SetMaxOpenConns(dbAccWrite.MaxOpenConns)
	d.DBRead.SetMaxOpenConns(dbAccRead.MaxOpenConns)

	d.DBRead.SetMaxIdleConns(dbAccRead.MaxIdleConns)
	d.DBWrite.SetMaxIdleConns(dbAccWrite.MaxIdleConns)
}

func (d *PostgreSQLHandler) Close() {
	if d.DBRead != nil {
		if err := d.DBRead.Close(); err != nil {
			log.Printf("Failed to close connection to DB Read! Error : %s\n", err.Error())
		} else {
			log.Println("Successfuly closing connection to DB Read")
		}
	}

	if d.DBWrite != nil {
		if err := d.DBWrite.Close(); err != nil {
			log.Printf("Failed to close connection to DB Write! Error : %s\n", err.Error())
		} else {
			log.Println("Successfuly closing connection to DB Write")
		}
	}
}

func (d *PostgreSQLHandler) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := d.DBWrite.Query(query, args...)
	return rows, err
}

func (d *PostgreSQLHandler) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := d.DBWrite.ExecContext(ctx, query, args...)
	return result, err
}

func (d *PostgreSQLHandler) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.DBWrite.Exec(query, args...)
	return result, err
}

func (d *PostgreSQLHandler) Get(dest interface{}, query string, args ...interface{}) error {
	err := d.DBRead.Get(dest, query, args...)
	return err
}

func (d *PostgreSQLHandler) GetMaster(dest interface{}, query string, args ...interface{}) error {
	err := d.DBWrite.Get(dest, query, args...)
	return err
}

func (d *PostgreSQLHandler) DriverName() string {
	return d.DBRead.DriverName()
}

func (d *PostgreSQLHandler) Select(dest interface{}, query string, args ...interface{}) error {
	err := d.DBRead.Select(dest, query, args...)
	return err
}

func (d *PostgreSQLHandler) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.DBRead.SelectContext(ctx, dest, query, args...)
}

func (d *PostgreSQLHandler) Begin() (IDBTx, error) {
	tx, err := d.DBWrite.Beginx()
	sqlxTx := SqlxTx{
		Tx: tx,
	}
	if err != nil {
		return &sqlxTx, err
	}
	return &sqlxTx, nil
}

func (d *PostgreSQLHandler) BeginTx() (*sqlx.Tx, error) {
	tx, err := d.DBWrite.Beginx()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (d *PostgreSQLHandler) BeginTxx(ctx context.Context) (*sqlx.Tx, error) {
	return d.DBWrite.BeginTxx(ctx, nil)
}

func (d *PostgreSQLHandler) Commit() error {
	return d.Tx.Commit()
}

func (d *PostgreSQLHandler) Rollback() error {
	return d.Tx.Rollback()
}

func (d *PostgreSQLHandler) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	tx, err := d.DBWrite.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (d *PostgreSQLHandler) TransactionBlock(tx *sqlx.Tx, fc func(tx *sqlx.Tx) error) error {
	if tx == nil {
		return errors.New("DB transaction is nil")
	}
	err := fc(tx)

	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return errTx
		}
		return err
	}

	errTx := tx.Commit()
	if errTx != nil {
		return errTx
	}

	return nil
}

func (d *PostgreSQLHandler) Rebind(query string) string {
	return d.DBRead.Rebind(query)
}

func (d *PostgreSQLHandler) In(query string, params ...interface{}) (string, []interface{}, error) {
	query, args, err := sqlx.In(query, params...)
	return query, args, err
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (d *PostgreSQLHandler) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.DBWrite.QueryRowContext(context.Background(), query, args...)
}

func (d *PostgreSQLHandler) QueryRowSqlx(query string, args ...interface{}) *sqlx.Row {
	return d.DBWrite.QueryRowx(query, args...)
}

func (d *PostgreSQLHandler) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := d.DBWrite.QueryContext(ctx, query, args...)
	return rows, err
}

func (d *PostgreSQLHandler) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := d.DBRead.GetContext(ctx, dest, query, args...)
	return err
}

func (d *PostgreSQLHandler) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	rows, err := d.DBWrite.QueryxContext(ctx, query, args...)
	return rows, err
}

func (d *PostgreSQLHandler) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return d.DBWrite.QueryRowxContext(ctx, query, args...)
}
