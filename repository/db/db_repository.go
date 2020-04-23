package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/domain/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
)

//Repository interface is a contract for all the methods contained by this db.Repository object.
type Repository interface {
	GetDishByID(int) (*dish.Dish, fcerr.FCErr)
	GetStorageByID(int) (*storage.Storage, fcerr.FCErr)
	GetDishes() (*dish.Dishes, fcerr.FCErr)
}

type repository struct {
	db *sql.DB
}

//NewRepository will get an instance of this type which satisfies the Repository interface.
func NewRepository(config string) (Repository, fcerr.FCErr) {
	db, err := sql.Open("mysql", strings.TrimSpace(config))
	if err != nil {
		fcerr := fcerr.NewInternalServerError("Error while connecting to the mysql database")
		return nil, fcerr
	}
	defer db.Close()

	//Check the connection to the database - If the credentials are wrong this will err out
	err = db.Ping()
	if err != nil {
		fcerr := fcerr.NewInternalServerError("Error while pinging the mysql database")
		return nil, fcerr
	}

	resultDB := repository{db}
	return &resultDB, nil
}

//GetDishes returns the list of all dishes in the database
func (repo *repository) GetDishes() (*dish.Dishes, fcerr.FCErr) {
	rows, err := repo.db.Query("Select * FROM dish")
	if err != nil {
		fcerr := fcerr.NewInternalServerError("Error while retrieving dishes from the database")
		return nil, fcerr
	}
	defer rows.Close()
	//s := "Retrieved Records:\n"

	for rows.Next() {
		var currentDish dish.Dish
		fmt.Println("There is a Next Row in the db! Before I scan, here is the variable currentDish:", currentDish)

		err := rows.Scan(&currentDish.DishID, &currentDish.UserID,
			&currentDish.StorageID, &currentDish.Title, &currentDish.Description,
			&currentDish.CreatedDate, &currentDish.ExpireDate, &currentDish.Priority)
		if err != nil {
			fcerr := fcerr.NewInternalServerError("unable to scan the result from the database")
			return nil, fcerr
		}

		fmt.Println("I just scanned, here is the variable currentDish:", currentDish)

	}

	return nil, nil
}

//GetDishByID takes an int and queries the mysql database for a dish with this id.
func (repo *repository) GetDishByID(int) (*dish.Dish, fcerr.FCErr) {
	var resultingDish dish.Dish
	return &resultingDish, nil
}

//GetStorageByID takes an int and queries the mysql database for a storage with this id.
func (repo *repository) GetStorageByID(int) (*storage.Storage, fcerr.FCErr) {
	var resultingStorage storage.Storage
	return &resultingStorage, nil
}
