package dish

import (
	"github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
	"github.com/jasonradcliffe/freshness-countdown-api/repository/db"
)

//Service is the interface that defines the contract for a dish service.
type Service interface {
	GetByID(int) (*dish.Dish, fcerr.FCErr)
	GetAll() (*dish.Dishes, fcerr.FCErr)
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
func (s *service) GetByID(id int) (*dish.Dish, fcerr.FCErr) {
	resultDish, err := s.repository.GetDishByID(id)
	if err != nil {
		return nil, fcerr.NewInternalServerError("could not do the GetByID, possibly not in the db")
	}
	return resultDish, nil
}

func (s *service) GetAll() (*dish.Dishes, fcerr.FCErr) {
	resultDishes, err := s.repository.GetDishes()
	if err != nil {
		fcerr := fcerr.NewInternalServerError("dish service could not do GetAll()")
		return nil, fcerr
	}
	return resultDishes, nil

}
