package storage

import (
	"fmt"

	dishDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/domain/storage"
	userDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/user"
	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
	"github.com/jasonradcliffe/freshness-countdown-api/repository/db"
)

//Service is the interface that defines the contract for a storage service.
type Service interface {
	GetByID(*userDomain.User, int) (*storage.Storage, fcerr.FCErr)
	GetDishesByID(*userDomain.User, int) (*dishDomain.Dishes, fcerr.FCErr)
	GetAll(*userDomain.User) (*storage.Storages, fcerr.FCErr)
	Create(*userDomain.User, *storage.Storage) (*storage.Storage, fcerr.FCErr)
	Update(*userDomain.User, *storage.Storage) fcerr.FCErr
	Delete(*userDomain.User, int) fcerr.FCErr
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

//GetByID: (alexaid string, accessToken string, id int) takes an int id and sends it to the database repo for lookup.
func (s *service) GetByID(requestingUser *userDomain.User, pID int) (*storage.Storage, fcerr.FCErr) {
	resultStorage, err := s.repository.GetStorageByID(requestingUser.UserID, pID)
	if err != nil {
		return nil, fcerr.NewInternalServerError("could not do the GetByID, possibly not in the db")
	}
	return resultStorage, nil
}

//GetDishesByID(requestingUser *userDomain.User, pID int) gets all the dishes that belong to the requesting user in the given storage unit
func (s *service) GetDishesByID(requestingUser *userDomain.User, pID int) (*dishDomain.Dishes, fcerr.FCErr) {
	resultDishes, err := s.repository.GetStorageDishes(requestingUser.UserID, pID)
	if err != nil {
		return nil, fcerr.NewInternalServerError("could not do the getstoragedishes")
	}
	return resultDishes, nil
}

//GetAll: (alexaid string, accessToken string) - gets all the storage units that the requesting user has
func (s *service) GetAll(requestUser *userDomain.User) (*storage.Storages, fcerr.FCErr) {
	resultStorageList, err := s.repository.GetStorages(requestUser.UserID)
	if err != nil {
		fcerr := fcerr.NewInternalServerError("storage service could not do GetAll()")
		return nil, fcerr
	}
	return resultStorageList, nil

}

func (s *service) Create(requestingUser *userDomain.User, newStorage *storage.Storage) (*storage.Storage, fcerr.FCErr) {

	personalCount, err := s.repository.GetPersonalStorageCount(requestingUser.UserID)
	if err != nil {
		return nil, fcerr.NewInternalServerError("Error when creating the storage unit.")
	}

	newStorage.UserID = requestingUser.UserID
	newStorage.PersonalID = personalCount + 1

	fmt.Println("\nWe are doing the storage service Create() with this storage:\n", newStorage)
	//alexaid string, accessToken string, storageID string, title string, desc string, expire string, priority string, dishtype string, portions string
	resultStorage, err := s.repository.CreateStorage(*newStorage)
	if err != nil {
		return nil, fcerr.NewInternalServerError("Storage Service could not do the Create()")
	}
	return resultStorage, nil

}

func (s *service) Update(requestingUser *userDomain.User, newStorage *storage.Storage) fcerr.FCErr {

	fmt.Println("\nWe are doing the storage service Update() with this storage:\n", newStorage)
	//alexaid string, accessToken string, storageID string, title string, desc string, expire string, priority string, dishtype string, portions string
	err := s.repository.UpdateStorage(*newStorage)
	if err != nil {
		return fcerr.NewInternalServerError("Storage Service could not do the Create()")
	}
	return nil
}

func (s *service) Delete(requestingUser *userDomain.User, storageID int) fcerr.FCErr {

	fmt.Println("We are doing the storage service Delete() with this storage:\n", storageID)
	//alexaid string, accessToken string, storageID string, title string, desc string, expire string, priority string, dishtype string, portions string
	err := s.repository.DeleteStorage(requestingUser.UserID, storageID)
	if err != nil {
		return fcerr.NewInternalServerError("Storage Service could not do the Delete()")
	}
	return nil

}
