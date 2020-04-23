package storage

import (
	"github.com/jasonradcliffe/freshness-countdown-api/domain/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
	"github.com/jasonradcliffe/freshness-countdown-api/repository/db"
)

//Service is the interface that defines the contract for a dish service.
type Service interface {
	GetByID(int) (*storage.Storage, fcerr.FCErr)
}

type service struct {
	repository db.Repository
}

//NewService takes a database repository and gives you a new Service instance.
func NewService(repo db.Repository) Service {
	return &service{
		repository: repo,
	}
}

//GetByID takes an int id and sends it to the database repo for lookup.
func (s *service) GetByID(id int) (*storage.Storage, fcerr.FCErr) {
	resultStorage, err := s.repository.GetStorageByID(id)
	if err != nil {
		return nil, fcerr.NewInternalServerError("could not do the GetByID for storage, possibly not in the db")
	}
	return resultStorage, nil
}
