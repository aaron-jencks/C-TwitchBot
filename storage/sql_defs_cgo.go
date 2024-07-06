//go:build cgo || arm.7
package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func getSqliteConn(fname string) (*sql.DB, error) {
	return sql.Open("sqlite3", fname)
}
