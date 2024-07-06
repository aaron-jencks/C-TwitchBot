//go:build !cgo
package storage

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func getSqliteConn(fname string) (*sql.DB, error) {
	return sql.Open("sqlite", fname)
}
