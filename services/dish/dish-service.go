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

//GetByID: (alexaid string, accessToken string, id int) takes an int id and sends it to the database repo for lookup.
func (s *service) GetByID(alexaid string, accessToken string, id int) (*dish.Dish, fcerr.FCErr) {
	resultDish, err := s.repository.GetDishByID(id)
	if err != nil {
		return nil, fcerr.NewInternalServerError("could not do the GetByID, possibly not in the db")
	}
	return resultDish, nil
}

//GetAll: (alexaid string, accessToken string) - gets all the dishes... if the user is admin
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

	fmt.Println("in the dish service Create(). Got this for expire window:\n" + dishMap["expireWindow"] + "\n")

	newDish := &dish.Dish{
		StorageID: newStorageID,
		Title:     dishMap["title"],
	}
	fmt.Println("\nWe are doing the dish service Create() with this dish:\n", newDish)
	//alexaid string, accessToken string, storageID string, title string, desc string, expire string, priority string, dishtype string, portions string
	resultDish, err := s.repository.CreateDish(*newDish)
	if err != nil {
		return nil, fcerr.NewInternalServerError("Dish Service could not do the Create()")
	}
	return resultDish, nil

}
