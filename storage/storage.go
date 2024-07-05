package storage

import "time"

type StorageBacking interface {
	// Counters
	CreateCounter(name string, initial int, prefix string) error
	RetrieveCounter(name string) (int, string, error)
	UpdateCounter(name string, newValue int) error
	DeleteCounter(name string) error
	ListCounters() ([]string, error)

	// Timers
	CreateTimer(name string, message string, interval time.Duration) error
	RetrieveTimer(name string) (string, time.Duration, time.Time, error)
	ResetTimer(name string) error
	DeleteTimer(name string) error
	ListTimers() (map[string]time.Time, error)
}
