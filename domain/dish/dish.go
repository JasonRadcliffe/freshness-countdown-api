package dish

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
)

//Dish type is the struct in the Domain that contains all the fields for what a Dish is.
type Dish struct {
	DishID         int    `json:"DishID"`
	PersonalDishID int    `json:"PersonalDishID"`
	UserID         int    `json:"UserID"`
	StorageID      int    `json:"StorageID"`
	Title          string `json:"Title"`
	Description    string `json:"Description"`
	CreatedDate    string `json:"TimeCreated"`
	ExpireDate     string `json:"TimeExpires"`
	Priority       string `json:"Priority"`
	DishType       string `json:"DishType"`
	Portions       int    `json:"Portions"`
	TempMatch      string `json:"TempMatch"`
}

//Dishes type is a slice of the domain type Dish.
type Dishes []Dish

//Contains methods and validators that a dish would know about
//isExpired()
//get new dish with title()
//change expiration()
//how long remaining()

//IsExpired will check the ExpireDate field against the current time, and return true for expired
func (d *Dish) IsExpired() (bool, fcerr.FCErr) {
	expireTime, err := time.Parse("2006-01-02T15:04", d.ExpireDate)
	if err != nil {
		fmt.Println("The dish did not have a valid expiration date. Error:", err.Error())
		return false, fcerr.NewInternalServerError("Encountered a dish without a valid expiration date")
	}

	if expireTime.After(time.Now()) {
		return false, nil
	}
	return true, nil
}

//WillExpireBy will check the ExpireDate field against the given date/time, and return true if the dish will be expired
func (d *Dish) WillExpireBy(dateStr string) (bool, fcerr.FCErr) {
	dishExpireTime, err := time.Parse("2006-01-02T15:04", d.ExpireDate)
	if err != nil {
		fmt.Println("The dish did not have a valid expiration date. Error:", err.Error())
		return false, fcerr.NewInternalServerError("Encountered a dish without a valid expiration date")
	}

	//checkTime, err := time.Parse("2006-01-02", dateStr)
	checkTime, err := dateparse.ParseAny(dateStr)
	if err != nil {
		fmt.Println("WillExpireBy was passed an invalid expiration string:" + dateStr)
		return false, fcerr.NewBadRequestError("dish method was passed an invalid expiration string")
	}

	if dishExpireTime.After(checkTime) {
		return false, nil
	}
	return true, nil
}
