package db

import (
	"database/sql"
	"fmt"

	"github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/domain/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
)

//Repository interface is a contract for all the methods contained by this db.Repository object.
type Repository interface {
	GetDishByID(int) (*dish.Dish, error)
	GetStorageByID(int) (*storage.Storage, error)
}

type repository struct {
	mysqldb *sql.DB
}

//NewRepository will get an instance of this type which satisfies the Repository interface.
func NewRepository(config string) (Repository, error) {
	fmt.Println("about to try to make a NewRepository with this config string:", config)
	db, err := sql.Open("mysql", config)
	if err != nil {
		fmt.Println("sql.Open() did in fact throw an error")
		fcerr := fcerr.NewInternalServerError("Error while connecting to the mysql database")
		return nil, fcerr
	}
	fmt.Println("sql.Open() did not error")
	defer db.Close()

	//Check the connection to the database - If the credentials are wrong this will err out
	fmt.Println("about to ping the connection to the database - if the credentials are wrong this will err out")
	err = db.Ping()
	if err != nil {
		fmt.Println("yep, got an error when pinging the db")
		fcerr := fcerr.NewInternalServerError("Error while pinging the mysql database")
		return nil, fcerr
	}

	fmt.Println("got all the way through to the other end of NewRepository()")

	resultDB := repository{db}
	return &resultDB, nil
}

//GetDishByID takes an int and queries the mysql database for a dish with this id.
func (r *repository) GetDishByID(int) (*dish.Dish, error) {
	var resultingDish dish.Dish
	return &resultingDish, nil
}

//GetStorageByID takes an int and queries the mysql database for a storage with this id.
func (r *repository) GetStorageByID(int) (*storage.Storage, error) {
	var resultingStorage storage.Storage
	return &resultingStorage, nil
}
