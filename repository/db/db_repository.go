package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/domain/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
)

const getDishesQuery = `SELECT * FROM dish`

const getDishByIDQuery = `SELECT * FROM dish WHERE id = %d`

const getDishByTempMatch = `Select * FROM dish WHERE temp_match = %s`

const createDishQuery = `INSERT INTO dish ` +
	`(user_id, storage_id, title, description, created_date, expire_date, priority, dish_type, portions, temp_match) ` +
	`VALUES(%d, %d, %s, %s, %s, %s, %s, %s, %d, %s);`

const updateDishQuery = `UPDATE dish SET storage_id = %s, title = %s, description = %s, expire_date = %s, ` +
	`priority = %s, dish_type = %s, portions = %d WHERE id=%d`

const deleteDishQuery = `DELETE FROM dish WHERE id=%d`

const getUsersQuery = `SELECT * FROM user`

const getUserByIDQuery = `SELECT * FROM user WHERE id = %d`

const getUserByEmailQuery = `SELECT * FROM user WHERE email = %s`

const createUserQuery = `INSERT INTO user (email, created_date, access_token, temp_match) ` +
	`VALUES(%s, %s, %s, %s)`

const deleteUserQuery = `DELETE FROM user WHERE id=%d`

const getAllStorageQuery = `SELECT * FROM storage`

const getStorageByIDQuery = `SELECT * FROM storage WHERE id=%d`

const createStorageQuery = `INSERT INTO storage (user_id, title, description, temp_match) ` +
	`VALUES(%d, %s, %s, %s)`

const updateStorageQuery = `UPDATE storage SET title = %s, description = %s WHERE id=%d`

const deleteStorageQuery = `DELETE FROM storage WHERE id=%d`

const getStorageDishesQuery = `SELECT * FROM dish WHERE storage_id = %d`

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

	//trying without the db.Close()
	//defer db.Close()

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
	fmt.Println("now at the beginning of the db_repository GetDishes()")
	var resultDishes dish.Dishes
	rows, err := repo.db.Query(`Select * FROM dish`)
	fmt.Println("now after doing the Query")
	if err != nil {
		fmt.Println("got an error on the Query")
		fcerr := fcerr.NewInternalServerError("Error while retrieving dishes from the database")
		return nil, fcerr
	}
	defer rows.Close()
	//s := "Retrieved Records:\n"
	fmt.Println("now about to check the rows returned:")
	for rows.Next() {
		var currentDish dish.Dish
		fmt.Println("Inside the result set loop. currentDish:", currentDish)
		err := rows.Scan(&currentDish.DishID, &currentDish.UserID, &currentDish.StorageID, &currentDish.Title,
			&currentDish.Description, &currentDish.CreatedDate, &currentDish.ExpireDate, &currentDish.Priority,
			&currentDish.DishType, &currentDish.Portions, &currentDish.TempMatch)
		if err != nil {
			fcerr := fcerr.NewInternalServerError("unable to scan the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current dish scanned. currentDish:", currentDish)
		resultDishes = append(resultDishes, currentDish)

	}

	return &resultDishes, nil
}

//GetDishByID takes an int and queries the mysql database for a dish with this id.
func (repo *repository) GetDishByID(int) (*dish.Dish, fcerr.FCErr) {
	var resultingDish dish.Dish
	return &resultingDish, nil
}

func (repo *repository) CreateDish() (*dish.Dish, fcerr.FCErr) {
	return nil, nil
}

//GetStorageByID takes an int and queries the mysql database for a storage with this id.
func (repo *repository) GetStorageByID(int) (*storage.Storage, fcerr.FCErr) {
	var resultingStorage storage.Storage
	return &resultingStorage, nil
}
