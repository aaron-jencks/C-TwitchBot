package storage

type StorageBacking interface {
	// Counters
	CreateCounter(name string, initial int, prefix string) error
	RetrieveCounter(name string) (int, error)
	UpdateCounter(name string, newValue int) error
	DeleteCounter(name string) error
}
