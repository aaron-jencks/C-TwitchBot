package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SqliteBackingStore struct {
	fname string
}

func (sb SqliteBackingStore) setupTables() error {
	db, err := sql.Open("sqlite", sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("create table if not exists counters (name text primary key, value integer, prefix text)")
	if err != nil {
		return err
	}
	return nil
}

func CreateSqliteBacker(fname string) (*SqliteBackingStore, error) {
	result := SqliteBackingStore{
		fname: fname,
	}
	err := result.setupTables()
	return &result, err
}

func (sb *SqliteBackingStore) CreateCounter(name string, initial int, prefix string) error {
	db, err := sql.Open("sqlite", sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("insert or replace into counters values (?, ?, ?)", name, initial, prefix)
	return err
}

func (sb *SqliteBackingStore) RetrieveCounter(name string) (value int, err error) {
	db, err := sql.Open("sqlite", sb.fname)
	if err != nil {
		return
	}
	defer db.Close()
	row := db.QueryRow("select value from counters where name = ?", name)
	err = row.Err()
	if err != nil {
		return
	}
	err = row.Scan(&value)
	return
}

func (sb *SqliteBackingStore) UpdateCounter(name string, newValue int) error {
	db, err := sql.Open("sqlite", sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("update or replace counters set value = ? where name = ?", newValue, name)
	return err
}

func (sb *SqliteBackingStore) DeleteCounter(name string) error {
	db, err := sql.Open("sqlite", sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("delete from counters where name = ?", name)
	return err
}
