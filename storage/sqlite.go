package storage

import (
	"database/sql"
	"time" 
)

type SqliteBackingStore struct {
	fname string
}

func (sb SqliteBackingStore) setupTables() error {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("create table if not exists counters (name text primary key, value integer, prefix text)")
	if err != nil {
		return err
	}
	_, err = db.Exec("create table if not exists timers (name text primary key, message text, interval integer, next text)")
	if err != nil {
		return err
	}
	_, err = db.Exec("create table if not exists mappings (name text primary key, message text)")
	if err != nil {
		return err
	}
	return nil
}

func (sb *SqliteBackingStore) GetDbConn() (*sql.DB, error) {
	return getSqliteConn(sb.fname)
}


func CreateSqliteBacker(fname string) (*SqliteBackingStore, error) {
	result := SqliteBackingStore{
		fname: fname,
	}
	err := result.setupTables()
	return &result, err
}

func (sb *SqliteBackingStore) CreateCounter(name string, initial int, prefix string) error {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("insert or ignore into counters values (?, ?, ?)", name, initial, prefix)
	return err
}

func (sb *SqliteBackingStore) RetrieveCounter(name string) (value int, prefix string, err error) {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return
	}
	defer db.Close()
	row := db.QueryRow("select value, prefix from counters where name = ?", name)
	err = row.Err()
	if err != nil {
		return
	}
	err = row.Scan(&value, &prefix)
	return
}

func (sb *SqliteBackingStore) UpdateCounter(name string, newValue int) error {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("update or replace counters set value = ? where name = ?", newValue, name)
	return err
}

func (sb *SqliteBackingStore) DeleteCounter(name string) error {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("delete from counters where name = ?", name)
	return err
}

func (sb *SqliteBackingStore) ListCounters() ([]string, error) {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select name from counters")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []string
	var temp string
	for rows.Next() {
		err = rows.Scan(&temp)
		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}
	err = rows.Err()
	return result, err
}

func (sb *SqliteBackingStore) CreateTimer(name string, message string, interval time.Duration) error {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	next := time.Now().Add(interval).Format(time.RFC3339)
	_, err = db.Exec("insert or ignore into timers values (?, ?, ?, ?)", name, message, interval.Nanoseconds(), next)
	return err
}

func (sb *SqliteBackingStore) RetrieveTimer(name string) (message string, interval time.Duration, next time.Time, err error) {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return
	}
	defer db.Close()
	row := db.QueryRow("select message, interval, next from timers where name = ?", name)
	err = row.Err()
	if err != nil {
		return
	}
	var iint int64
	var snext string
	err = row.Scan(&message, &iint, &snext)
	if err != nil {
		return
	}
	interval = time.Duration(iint)
	next, err = time.Parse(time.RFC3339, snext)
	return
}

func (sb *SqliteBackingStore) ResetTimer(name string) error {
	_, intr, _, err := sb.RetrieveTimer(name)
	if err != nil {
		return err
	}

	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	next := time.Now().Add(intr).Format(time.RFC3339)
	_, err = db.Exec("update or replace timers set next = ? where name = ?", name, next)
	return err
}

func (sb *SqliteBackingStore) DeleteTimer(name string) error {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("delete from timers where name = ?", name)
	return err
}

func (sb *SqliteBackingStore) ListTimers() (map[string]time.Time, error) {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select name, next from timers order by strftime(\"%s\", next)")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]time.Time{}
	var stemp string
	var ntemp string
	for rows.Next() {
		err = rows.Scan(&stemp, &ntemp)
		if err != nil {
			return nil, err
		}
		result[stemp], err = time.Parse(time.RFC3339, ntemp)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	return result, err
}

func (sb *SqliteBackingStore) CreateMapping(name, message string) error {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("insert or ignore into mappings values (?, ?)", name, message)
	return err
}

func (sb *SqliteBackingStore) RetrieveMapping(name string) (msg string, err error) {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return
	}
	defer db.Close()
	row := db.QueryRow("select message from mappings where name = ?", name)
	err = row.Err()
	if err != nil {
		return
	}
	err = row.Scan(&msg)
	return
}

func (sb *SqliteBackingStore) UpdateMapping(name, newMessage string) error {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("update or replace mappings set message = ? where name = ?", newMessage, name)
	return err
}

func (sb *SqliteBackingStore) DeleteMapping(name string) error {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("delete from mappings where name = ?", name)
	return err
}

func (sb *SqliteBackingStore) ListMappings() (result map[string]string, err error) {
	db, err := getSqliteConn(sb.fname)
	if err != nil {
		return
	}
	defer db.Close()
	rows, err := db.Query("select * from mappings order by name")
	if err != nil {
		return
	}
	defer rows.Close()
	result = map[string]string{}
	var tname, tmsg string
	for rows.Next() {
		err = rows.Scan(&tname, &tmsg)
		if err != nil {
			return
		}
		result[tname] = tmsg
	}
	err = rows.Err()
	return
}
