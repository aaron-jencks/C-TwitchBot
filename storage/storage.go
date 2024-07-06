package storage

import (
	"database/sql"
	"time"
)

type StorageBacking interface {
	// General
	GetDbConn() (*sql.DB, error)

	// Counters
	CreateCounter(name string, initial int, prefix string) error
	RetrieveCounter(name string) (int, string, error)
	UpdateCounter(name string, newValue int) error
	DeleteCounter(name string) error
	ListCounters() ([]string, error)

	// Timers
	CreateTimer(name, message string, interval time.Duration) error
	RetrieveTimer(name string) (string, time.Duration, time.Time, error)
	ResetTimer(name string) error
	DeleteTimer(name string) error
	ListTimers() (map[string]time.Time, error)

	// Mappings
	CreateMapping(name, message string) error
	RetrieveMapping(name string) (string, error)
	UpdateMapping(name, newMessage string) error
	DeleteMapping(name string) error
	ListMappings() (map[string]string, error)
}
