package dish

import (
	"fmt"
	"strconv"

	"github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
	"github.com/jasonradcliffe/freshness-countdown-api/repository/db"
)

//Service is the interface that defines the contract for a dish service.
type Service interface {
	GetByID(string, string, int) (*dish.Dish, fcerr.FCErr)
	GetAll(string, string) (*dish.Dishes, fcerr.FCErr)
	Create(string, string, map[string]string) (*dish.Dish, fcerr.FCErr)
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
func (s *service) GetByID(alexaid string, accessToken string, id int) (*dish.Dish, fcerr.FCErr) {
	resultDish, err := s.repository.GetDishByID(id)
	if err != nil {
		return nil, fcerr.NewInternalServerError("could not do the GetByID, possibly not in the db")
	}
	return resultDish, nil
}

func (s *service) GetAll(alexaid string, accessToken string) (*dish.Dishes, fcerr.FCErr) {
	resultDishes, err := s.repository.GetDishes()
	if err != nil {
		fcerr := fcerr.NewInternalServerError("dish service could not do GetAll()")
		return nil, fcerr
	}
	return resultDishes, nil

}

func (s *service) Create(alexaid string, accessToken string, dishMap map[string]string) (*dish.Dish, fcerr.FCErr) {
	newStorageID, err := strconv.Atoi(dishMap["storageID"])
	if err != nil {
		return nil, fcerr.NewBadRequestError("storage id was not a number")
	}

	newPortions, err := strconv.Atoi(dishMap["portions"])
	if err != nil {
		return nil, fcerr.NewBadRequestError("portions was not a number")
	}
	newDish := &dish.Dish{
		StorageID:   newStorageID,
		Title:       dishMap["title"],
		Description: dishMap["description"],
		Priority:    dishMap["priority"],
		DishType:    dishMap["dishType"],
		Portions:    newPortions,
	}
	fmt.Println("\nWe are doing the dish service Create() with this dish:\n", newDish)
	//alexaid string, accessToken string, storageID string, title string, desc string, expire string, priority string, dishtype string, portions string
	resultDish, err := s.repository.CreateDish(*newDish)
	if err != nil {
		return nil, fcerr.NewInternalServerError("Dish Service could not do the Create()")
	}
	return resultDish, nil

}
