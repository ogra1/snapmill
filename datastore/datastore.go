package datastore

import "github.com/ogra1/fabrica/domain"

// Datastore interface for the database logic
type Datastore interface {
	BuildList() ([]domain.Build, error)
	BuildCreate(name, repo string) (string, error)
	BuildUpdate(id, status string) error
	BuildGet(id string) (domain.Build, error)
	BuildLogCreate(id, message string) error
	BuildLogList(id string) ([]domain.BuildLog, error)
}