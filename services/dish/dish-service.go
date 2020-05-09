package dish

import (
	"fmt"
	"time"

	"github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	userDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/user"
	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
	"github.com/jasonradcliffe/freshness-countdown-api/repository/db"
)

//Service is the interface that defines the contract for a dish service.
type Service interface {
	GetByID(*userDomain.User, int) (*dish.Dish, fcerr.FCErr)
	GetExpired(*userDomain.User) (*dish.Dishes, fcerr.FCErr)
	GetAll(*userDomain.User) (*dish.Dishes, fcerr.FCErr)
	Create(*userDomain.User, *dish.Dish, string) (*dish.Dish, fcerr.FCErr)
	Update(*userDomain.User, *dish.Dish, string) fcerr.FCErr
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
func (s *service) GetByID(requestingUser *userDomain.User, pID int) (*dish.Dish, fcerr.FCErr) {
	resultDish, err := s.repository.GetDishByID(requestingUser.UserID, pID)
	if err != nil {
		return nil, fcerr.NewInternalServerError("could not do the GetByID, possibly not in the db")
	}
	return resultDish, nil
}

//GetAll: (alexaid string, accessToken string) - gets all the dishes... if the user is admin
func (s *service) GetAll(requestUser *userDomain.User) (*dish.Dishes, fcerr.FCErr) {
	resultDishes, err := s.repository.GetDishes(requestUser.UserID)
	if err != nil {
		fcerr := fcerr.NewInternalServerError("dish service could not do GetAll()")
		return nil, fcerr
	}
	return resultDishes, nil

}

//GetAll: (alexaid string, accessToken string) - gets all the dishes... if the user is admin
func (s *service) GetExpired(requestUser *userDomain.User) (*dish.Dishes, fcerr.FCErr) {
	//var cDish dish.Dish
	var expiredDishes dish.Dishes
	resultDishes, err := s.repository.GetDishes(requestUser.UserID)

	if err != nil {
		return nil, fcerr.NewInternalServerError("Could not retrieve the dishes")
	}

	for i, d := range *resultDishes {
		fmt.Println(i, "In the for each loop of the GetExpired!! dish Expire date:", d.ExpireDate)
		check, err := d.IsExpired()
		if err != nil {
			continue
		}
		if check == true {
			fmt.Println("Got a true - an expired dish!", d.Title, d.ExpireDate)
			expiredDishes = append(expiredDishes, d)
		}

	}

	return &expiredDishes, nil

}

func (s *service) Create(requestingUser *userDomain.User, newDish *dish.Dish, expireWindow string) (*dish.Dish, fcerr.FCErr) {

	//TODO: write conversions between Alexa duration and time.Now
	expireDate := "2020-10-13T08:00"
	datePattern := "2006-01-02T15:04:05"
	timeNow := time.Now().In(time.UTC)
	createdDate := timeNow.Format(datePattern)

	personalCount, err := s.repository.GetPersonalDishCount(*requestingUser)
	if err != nil {
		return nil, fcerr.NewInternalServerError("Error when creating the dish.")
	}

	newDish.UserID = requestingUser.UserID
	newDish.PersonalDishID = personalCount + 1
	newDish.CreatedDate = createdDate
	newDish.ExpireDate = expireDate

	fmt.Println("\nWe are doing the dish service Create() with this dish:\n", newDish)
	//alexaid string, accessToken string, storageID string, title string, desc string, expire string, priority string, dishtype string, portions string
	resultDish, err := s.repository.CreateDish(*newDish)
	if err != nil {
		return nil, fcerr.NewInternalServerError("Dish Service could not do the Create()")
	}
	return resultDish, nil

}

func (s *service) Update(requestingUser *userDomain.User, newDish *dish.Dish, expireWindow string) fcerr.FCErr {

	//TODO: write conversions between Alexa duration and time.Now
	expireDate := "2020-10-13T08:00"
	//datePattern := "2006-01-02T15:04"

	newDish.ExpireDate = expireDate

	fmt.Println("\nWe are doing the dish service Update() with this dish:\n", newDish)
	//alexaid string, accessToken string, storageID string, title string, desc string, expire string, priority string, dishtype string, portions string
	err := s.repository.UpdateDish(*newDish)
	if err != nil {
		return fcerr.NewInternalServerError("Dish Service could not do the Create()")
	}
	return nil

}
