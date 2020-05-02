package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/domain/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/domain/user"
	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
)

//GetDishesBase is the Query for GetDishes().
const GetDishesBase = `SELECT * FROM dish`

//GetDishByIDBase can be used with fmt.Sprintf() to get the Query for GetDishByID().
const GetDishByIDBase = `SELECT * FROM dish WHERE id = %d`

//GetDishByTempMatchBase can be used with fmt.Sprintf() to get the Query for GetDishByTempMatch().
const GetDishByTempMatchBase = `Select * FROM dish WHERE temp_match = "%s"`

//CreateDishBase can be used with fmt.Sprintf() to get the Query for CreateDish().
const CreateDishBase = `INSERT INTO dish ` +
	`(user_id, storage_id, title, description, created_date, expire_date, priority, dish_type, portions, temp_match) ` +
	`VALUES(%d, %d, "%s", "%s", "%s", "%s", "%s", "%s", %d, "%s")`

//UpdateDishBase can be used with fmt.Sprintf() to get the Query for UpdateDish().
const UpdateDishBase = `UPDATE dish SET storage_id = "%s", title = "%s", description = "%s", expire_date = "%s", ` +
	`priority = "%s", dish_type = "%s", portions = %d WHERE id=%d`

//DeleteDishBase can be used with fmt.Sprintf() to get the Query for DeleteDish().
const DeleteDishBase = `DELETE FROM dish WHERE id=%d`

//GetUsersBase is the Query for GetUsers().
const GetUsersBase = `SELECT * FROM user`

//GetUserByIDBase can be used with fmt.Sprintf() to get the Query for GetUserByID().
const GetUserByIDBase = `SELECT * FROM user WHERE id = %d`

//GetUserByEmailBase can be used with fmt.Sprintf() to get the Query for GetUserByEmail().
const GetUserByEmailBase = `SELECT * FROM user WHERE email = "%s"`

//GetUserByAlexaBase can be used with fmt.Sprintf() to get the Query for GetUserByAlexa().
const GetUserByAlexaBase = `SELECT * FROM user WHERE alexa_user_id = "%s"`

//GetUserByTempMatchBase can be used with fmt.Sprintf() to get the Query for GetUserByTempMatch().
const GetUserByTempMatchBase = `SELECT * FROM user WHERE temp_match = "%s"`

//CreateUserBase can be used with fmt.Sprintf() to get the Query for CreateUser().
const CreateUserBase = `INSERT INTO user (email, first_name, last_name, full_name, created_date, access_token, refresh_token, alexa_user_id, temp_match) ` +
	`VALUES("%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s")`

//UpdateUserBase can be used with fmt.Sprintf() to get the Query for UpdateUser().
const UpdateUserBase = `UPDATE user SET email = "%s", first_name = "%s", last_name = "%s", full_name = "%s", ` +
	`created_date = "%s", access_token = "%s", refresh_token = "%s", alexa_user_id = "%s", temp_match = "%s" `

//DeleteUserBase can be used with fmt.Sprintf() to get the Query for DeleteUser().
const DeleteUserBase = `DELETE FROM user WHERE id=%d`

//GetAllStorageBase can be used with fmt.Sprintf() to get the Query for GetAllStorage().
const GetAllStorageBase = `SELECT * FROM storage`

//GetStorageByIDBase can be used with fmt.Sprintf() to get the Query for GetStorageByID().
const GetStorageByIDBase = `SELECT * FROM storage WHERE id=%d`

//CreateStorageBase can be used with fmt.Sprintf() to get the Query for CreateStorage().
const CreateStorageBase = `INSERT INTO storage (user_id, title, description, temp_match) ` +
	`VALUES(%d, "%s", "%s", "%s")`

//UpdateStorageBase can be used with fmt.Sprintf() to get the Query for UpdateStorage().
const UpdateStorageBase = `UPDATE storage SET title = "%s", description = "%s" WHERE id=%d`

//DeleteStorageBase can be used with fmt.Sprintf() to get the Query for DeleteStorage().
const DeleteStorageBase = `DELETE FROM storage WHERE id=%d`

//GetStorageDishesBase can be used with fmt.Sprintf() to get the Query for GetStorageDishes().
const GetStorageDishesBase = `SELECT * FROM dish WHERE storage_id = %d`

//Repository interface is a contract for all the methods contained by this db.Repository object.
type Repository interface {
	GetDishes() (*dish.Dishes, fcerr.FCErr)
	GetDishByID(int) (*dish.Dish, fcerr.FCErr)
	GetDishByTempMatch(string) (*dish.Dish, fcerr.FCErr)
	CreateDish(dish.Dish) (*dish.Dish, fcerr.FCErr)
	UpdateDish(dish.Dish) (*dish.Dish, fcerr.FCErr)
	DeleteDish(dish.Dish) fcerr.FCErr

	GetUsers() (*user.Users, fcerr.FCErr)
	GetUserByID(int) (*user.User, fcerr.FCErr)
	GetUserByEmail(string) (*user.User, fcerr.FCErr)
	GetUserByAlexa(string) (*user.User, fcerr.FCErr)
	GetUserByTempMatch(string) (*user.User, fcerr.FCErr)
	CreateUser(user.User) (*user.User, fcerr.FCErr)
	UpdateUser(user.User) (*user.User, fcerr.FCErr)
	DeleteUser(user.User) fcerr.FCErr

	GetStorage(int) (*storage.Storages, fcerr.FCErr)
	GetStorageByID(int) (*storage.Storage, fcerr.FCErr)
	CreateStorage(storage.Storage) (*storage.Storage, fcerr.FCErr)
	UpdateStorage(storage.Storage) (*storage.Storage, fcerr.FCErr)
	DeleteStorage(storage.Storage) fcerr.FCErr

	GetStorageDishes(int) (*dish.Dishes, fcerr.FCErr)
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
	getDishesQuery := fmt.Sprintf(GetDishesBase)
	rows, err := repo.db.Query(getDishesQuery)
	fmt.Println("now after doing the Query:", getDishesQuery)
	if err != nil {
		fmt.Println("got an error on the Query:", err.Error())
		fcerr := fcerr.NewInternalServerError("Error while retrieving dishes from the database")
		return nil, fcerr
	}
	defer rows.Close()
	//s := "Retrieved Records:\n"
	fmt.Println("now about to check the rows returned:")
	count := 0
	for rows.Next() {
		count++
		var currentDish dish.Dish
		fmt.Println("Inside the result set loop. currentDish:", currentDish)
		err := rows.Scan(&currentDish.DishID, &currentDish.UserID, &currentDish.StorageID, &currentDish.Title,
			&currentDish.Description, &currentDish.CreatedDate, &currentDish.ExpireDate, &currentDish.Priority,
			&currentDish.DishType, &currentDish.Portions, &currentDish.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fmt.Println("&currentDish.DishID:", currentDish.DishID)
			fmt.Println("&currentDish.TempMatch:", currentDish.TempMatch)
			fcerr := fcerr.NewInternalServerError("Error while scanning the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current dish scanned. currentDish:", currentDish)
		resultDishes = append(resultDishes, currentDish)

	}
	if count < 1 {
		fcerr := fcerr.NewNotFoundError("Database could not find any dishes")
		fmt.Println("Database could not find any dishes")
		return nil, fcerr
	}

	return &resultDishes, nil
}

//GetDishByID takes an int and queries the mysql database for a dish with this id.
func (repo *repository) GetDishByID(id int) (*dish.Dish, fcerr.FCErr) {
	var resultingDish dish.Dish
	getDishByIDQuery := fmt.Sprintf(GetDishByIDBase, id)
	fmt.Println("about to run this query in GetDishByID:", getDishByIDQuery)

	rows, err := repo.db.Query(getDishByIDQuery)
	fmt.Println("now after doing the Query:", getDishByIDQuery)
	if err != nil {
		fmt.Println("got an error on the Query:", err.Error())
		fcerr := fcerr.NewInternalServerError("Error while retrieving dish from the database")
		return nil, fcerr
	}
	defer rows.Close()
	//s := "Retrieved Records:\n"
	fmt.Println("now about to check the rows returned:")
	var count = 0
	for rows.Next() {
		count++
		if count > 1 {
			dberr := fcerr.NewInternalServerError("Database returned more than 1 row when only 1 was expected")
			return nil, dberr
		}

		var currentDish dish.Dish
		fmt.Println("Inside the result set loop. currentDish:", currentDish)
		err := rows.Scan(&currentDish.DishID, &currentDish.UserID, &currentDish.StorageID, &currentDish.Title,
			&currentDish.Description, &currentDish.CreatedDate, &currentDish.ExpireDate, &currentDish.Priority,
			&currentDish.DishType, &currentDish.Portions, &currentDish.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fmt.Println("&currentDish.DishID:", currentDish.DishID)
			fmt.Println("&currentDish.TempMatch:", currentDish.TempMatch)
			fcerr := fcerr.NewInternalServerError("Error while scanning the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current dish scanned. currentDish:", currentDish)
		resultingDish = currentDish

	}
	if count == 0 {
		fcerr := fcerr.NewNotFoundError("Database could not find a dish with this ID")
		return nil, fcerr
	}
	return &resultingDish, nil

}

//GetDishByTempMatch takes a string and queries the mysql database for a dish with this temp_match.
func (repo *repository) GetDishByTempMatch(tm string) (*dish.Dish, fcerr.FCErr) {
	var resultingDish dish.Dish
	getDishByTempMatchQuery := fmt.Sprintf(GetDishByTempMatchBase, tm)
	fmt.Println("about to run this query in GetDishByTempMatch:", getDishByTempMatchQuery)

	rows, err := repo.db.Query(getDishByTempMatchQuery)
	fmt.Println("now after doing the Query:", getDishByTempMatchQuery)
	if err != nil {
		fmt.Println("got an error on the Query:", err.Error())
		fcerr := fcerr.NewInternalServerError("Error while retrieving dish from the database")
		return nil, fcerr
	}
	defer rows.Close()
	//s := "Retrieved Records:\n"
	fmt.Println("now about to check the rows returned:")
	var count = 0
	for rows.Next() {
		count++
		if count > 1 {
			dberr := fcerr.NewInternalServerError("Database returned more than 1 row when only 1 was expected")
			return nil, dberr
		}

		var currentDish dish.Dish
		fmt.Println("Inside the result set loop. currentDish:", currentDish)
		err := rows.Scan(&currentDish.DishID, &currentDish.UserID, &currentDish.StorageID, &currentDish.Title,
			&currentDish.Description, &currentDish.CreatedDate, &currentDish.ExpireDate, &currentDish.Priority,
			&currentDish.DishType, &currentDish.Portions, &currentDish.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fmt.Println("&currentDish.DishID:", currentDish.DishID)
			fmt.Println("&currentDish.TempMatch:", currentDish.TempMatch)
			fcerr := fcerr.NewInternalServerError("Error while scanning the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current dish scanned. currentDish:", currentDish)
		resultingDish = currentDish

	}
	if count == 0 {
		fcerr := fcerr.NewNotFoundError("Database could not find a dish with this temp match")
		return nil, fcerr
	}
	return &resultingDish, nil
}

//CreateDish takes a dish object and tries to add it to the database
func (repo *repository) CreateDish(d dish.Dish) (*dish.Dish, fcerr.FCErr) {
	createDishQuery := fmt.Sprintf(CreateDishBase, d.UserID, d.StorageID, d.Title, d.Description,
		d.CreatedDate, d.ExpireDate, d.Priority, d.DishType, d.Portions, d.TempMatch)

	fmt.Println("About to run this Query on the database:\n", createDishQuery)

	_, err := repo.db.Query(createDishQuery)
	if err != nil {
		fmt.Println("got an error on the Query:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while inserting the dish into the database")
		return nil, fcerr
	}

	checkDish, err := repo.GetDishByTempMatch(d.TempMatch)
	if err != nil {
		fmt.Println("Trying to CreateDish, seem to have hit a snag. Got an error when checking what we just put in: " + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while checking the dish that was created." +
			" Cannot verify if anything was entered to the Database")
		return nil, fcerr
	}

	return checkDish, nil
}

//UpdateDish takes a dish object and tries to update the existing dish in the database to match
func (repo *repository) UpdateDish(d dish.Dish) (*dish.Dish, fcerr.FCErr) {
	return nil, nil
}

//DeleteDish takes a dish object and tries to delete the existing dish from the database
func (repo *repository) DeleteDish(d dish.Dish) fcerr.FCErr {
	return nil
}

//GetUsers queries the database and returns a slice of User objects
func (repo *repository) GetUsers() (*user.Users, fcerr.FCErr) {
	return nil, nil
}

//GetUserByID gets a user from the database with the given ID.
func (repo *repository) GetUserByID(id int) (*user.User, fcerr.FCErr) {
	getUserByIDQuery := fmt.Sprintf(GetUserByIDBase, id)
	fmt.Println("About to run this Query on the database:\n", getUserByIDQuery)
	var resultingUser user.User

	rows, err := repo.db.Query(getUserByIDQuery)
	if err != nil {
		fmt.Println("got an error on the Query")
		fcerr := fcerr.NewInternalServerError("Error while retrieving user from the database by id")
		return nil, fcerr
	}
	defer rows.Close()
	//s := "Retrieved Records:\n"
	fmt.Println("now about to check the rows returned:")
	for rows.Next() {
		var cUser user.User
		fmt.Println("Inside the result set loop. currentUser:", cUser)
		err := rows.Scan(&cUser.UserID, &cUser.Email, &cUser.FirstName, &cUser.LastName,
			&cUser.FullName, &cUser.CreatedDate, &cUser.AccessToken, &cUser.RefreshToken, &cUser.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fcerr := fcerr.NewInternalServerError("unable to scan the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current user scanned. currentUser:", cUser)
		resultingUser = cUser

	}

	return &resultingUser, nil
}

//GetUserByEmail gets a user from the database with the given email.
func (repo *repository) GetUserByEmail(email string) (*user.User, fcerr.FCErr) {
	getUserByEmailQuery := fmt.Sprintf(GetUserByEmailBase, email)
	fmt.Println("About to run this Query on the database:\n", getUserByEmailQuery)
	var resultingUser user.User

	rows, err := repo.db.Query(getUserByEmailQuery)
	if err != nil {
		fmt.Println("got an error on the Query")
		fcerr := fcerr.NewInternalServerError("Error while retrieving user from the database by email")
		return nil, fcerr
	}
	defer rows.Close()
	//s := "Retrieved Records:\n"
	fmt.Println("now about to check the rows returned:")
	for rows.Next() {
		var cUser user.User
		fmt.Println("Inside the result set loop. currentUser:", cUser)
		err := rows.Scan(&cUser.UserID, &cUser.Email, &cUser.FirstName, &cUser.LastName,
			&cUser.FullName, &cUser.CreatedDate, &cUser.AccessToken, &cUser.RefreshToken, &cUser.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fcerr := fcerr.NewInternalServerError("unable to scan the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current user scanned. currentUser:", cUser)
		resultingUser = cUser

	}

	return &resultingUser, nil
}

//GetUserByAlexa gets a user from the database with the given alexa_user_id.
func (repo *repository) GetUserByAlexa(aID string) (*user.User, fcerr.FCErr) {
	getUserByAlexaQuery := fmt.Sprintf(GetUserByAlexaBase, aID)
	fmt.Println("About to run this Query on the database:\n", getUserByAlexaQuery)
	var resultingUser user.User

	rows, err := repo.db.Query(getUserByAlexaQuery)
	if err != nil {
		fmt.Println("got an error on the Query")
		fcerr := fcerr.NewInternalServerError("Error while retrieving user from the database by email")
		return nil, fcerr
	}
	defer rows.Close()
	//s := "Retrieved Records:\n"
	fmt.Println("now about to check the rows returned:")
	for rows.Next() {
		var cUser user.User
		fmt.Println("Inside the result set loop. currentUser:", cUser)
		err := rows.Scan(&cUser.UserID, &cUser.Email, &cUser.FirstName, &cUser.LastName,
			&cUser.FullName, &cUser.CreatedDate, &cUser.AccessToken, &cUser.RefreshToken, &cUser.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fcerr := fcerr.NewInternalServerError("unable to scan the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current user scanned. currentUser:", cUser)
		resultingUser = cUser

	}

	return &resultingUser, nil
}

//GetUserByTempMatch gets a user from the database with the given email.
func (repo *repository) GetUserByTempMatch(tm string) (*user.User, fcerr.FCErr) {
	getUserByTempMatchQuery := fmt.Sprintf(GetUserByTempMatchBase, tm)
	fmt.Println("About to run this Query on the database:\n", getUserByTempMatchQuery)
	var resultingUser user.User

	rows, err := repo.db.Query(getUserByTempMatchQuery)
	if err != nil {
		fmt.Println("got an error on the Query")
		fcerr := fcerr.NewInternalServerError("Error while retrieving user from the database by temp match")
		return nil, fcerr
	}
	defer rows.Close()
	//s := "Retrieved Records:\n"
	fmt.Println("now about to check the rows returned:")
	for rows.Next() {
		var cUser user.User
		fmt.Println("Inside the result set loop. currentUser:", cUser)
		err := rows.Scan(&cUser.UserID, &cUser.Email, &cUser.FirstName, &cUser.LastName,
			&cUser.FullName, &cUser.CreatedDate, &cUser.AccessToken, &cUser.RefreshToken, &cUser.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fcerr := fcerr.NewInternalServerError("unable to scan the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current user scanned. currentUser:", cUser)
		resultingUser = cUser

	}

	return &resultingUser, nil
}

//CreateUser adds a user to the database after being populated by the service.
func (repo *repository) CreateUser(u user.User) (*user.User, fcerr.FCErr) {
	createUserQuery := fmt.Sprintf(CreateUserBase, u.Email, u.FirstName, u.LastName, u.FullName, u.CreatedDate, u.AccessToken, u.RefreshToken, u.AlexaUserID, u.TempMatch)
	fmt.Println("About to run this Query on the database:\n", createUserQuery)

	_, err := repo.db.Query(createUserQuery)
	if err != nil {
		fmt.Println("got an error on the Query")
		fcerr := fcerr.NewInternalServerError("Error while inserting the user into the database")
		return nil, fcerr
	}

	checkUser, err := repo.GetUserByTempMatch(u.TempMatch)
	if err != nil {
		fmt.Println("Trying to CreateUser, seem to have hit a snag. Got an error when checking what we just put in:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while checking the user that was created")
		return nil, fcerr
	}

	return checkUser, nil
}

//UpdateUser takes a user object and tries to update the existing user in the database to match
func (repo *repository) UpdateUser(u user.User) (*user.User, fcerr.FCErr) {
	return nil, nil
}

//DeleteUser takes a user object and tries to delete the existing user from the database
func (repo *repository) DeleteUser(u user.User) fcerr.FCErr {
	return nil
}

//GetStorage takes an int of a user id and returns the list of storage objects owned by that user.
func (repo *repository) GetStorage(userID int) (*storage.Storages, fcerr.FCErr) {
	var resultingStorages storage.Storages
	return &resultingStorages, nil
}

//GetStorageByID takes an int and queries the mysql database for a storage with this id.
func (repo *repository) GetStorageByID(id int) (*storage.Storage, fcerr.FCErr) {
	var resultingStorage storage.Storage
	return &resultingStorage, nil
}

//CreateStorage takes a storage object and tries to add it to the database
func (repo *repository) CreateStorage(s storage.Storage) (*storage.Storage, fcerr.FCErr) {
	return nil, nil
}

//UpdateStorage takes a storage object and tries to update the existing storage in the database to match
func (repo *repository) UpdateStorage(s storage.Storage) (*storage.Storage, fcerr.FCErr) {
	return nil, nil
}

//DeleteStorage takes a storage object and tries to delete the existing storage from the database
func (repo *repository) DeleteStorage(s storage.Storage) fcerr.FCErr {
	return nil
}

//GetStorageDishes takes a storage object and tries to update the existing storage in the database to match
func (repo *repository) GetStorageDishes(s int) (*dish.Dishes, fcerr.FCErr) {
	return nil, nil
}
