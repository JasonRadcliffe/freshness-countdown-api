package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/domain/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/domain/user"
	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
)

//GetDishesBase is the Query for GetDishes().
const GetDishesBase = `SELECT * FROM dish WHERE user_id = %d`

//GetDishByIDBase can be used with fmt.Sprintf() to get the Query for GetDishByID().
const GetDishByIDBase = `SELECT * FROM dish WHERE user_id = %d AND personal_id = %d`

//GetDishByTempMatchBase can be used with fmt.Sprintf() to get the Query for GetDishByTempMatch().
const GetDishByTempMatchBase = `SELECT * FROM dish WHERE temp_match = "%s"`

//GetPersonalDishCountBase returns the number of dishes a given user has in the database, to be used for personal_id field
const GetPersonalDishCountBase = `SELECT COUNT(*) FROM dish WHERE user_id = %d`

//DecrementSomeDishesBase is used to shift every dish "up" after one in the middle of the dish list is deleted
const DecrementSomeDishesBase = `UPDATE dish SET personal_id = personal_id - 1 WHERE user_id = %d AND personal_id IN(%s)`

//CreateDishBase can be used with fmt.Sprintf() to get the Query for CreateDish().
const CreateDishBase = `INSERT INTO dish ` +
	`(personal_id, user_id, storage_id, title, description, created_date, expire_date, priority, dish_type, portions, temp_match) ` +
	`VALUES(%d, %d, %d, "%s", "%s", "%s", "%s", "%s", "%s", %d, "%s")`

//UpdateDishBase can be used with fmt.Sprintf() to get the Query for UpdateDish().
const UpdateDishBase = `UPDATE dish SET personal_id = %d, storage_id = %d, title = "%s", description = "%s", expire_date = "%s", ` +
	`priority = "%s", dish_type = "%s", portions = %d WHERE id=%d`

//DeleteDishBase can be used with fmt.Sprintf() to get the Query for DeleteDish().
const DeleteDishBase = `DELETE FROM dish WHERE user_id = %d AND personal_id=%d`

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
const CreateUserBase = `INSERT INTO user (email, first_name, last_name, full_name, created_date, access_token, refresh_token, alexa_user_id, is_admin, temp_match) ` +
	`VALUES("%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", %t, "%s")`

//UpdateUserBase can be used with fmt.Sprintf() to get the Query for UpdateUser().
const UpdateUserBase = `UPDATE user SET email = "%s", first_name = "%s", last_name = "%s", full_name = "%s", ` +
	`access_token = "%s", refresh_token = "%s", alexa_user_id = "%s", temp_match = "%s" WHERE id = %d `

//DeleteUserBase can be used with fmt.Sprintf() to get the Query for DeleteUser().
const DeleteUserBase = `DELETE FROM user WHERE id=%d`

//GetStoragesBase can be used with fmt.Sprintf() to get the Query for GetAllStorage().
const GetStoragesBase = `SELECT * FROM storage WHERE user_id=%d`

//GetStorageByIDBase can be used with fmt.Sprintf() to get the Query for GetStorageByID().
const GetStorageByIDBase = `SELECT * FROM storage WHERE user_id = %d AND personal_id = %d`

//GetStorageByTempMatchBase can be used with fmt.Sprintf() to get the Query for GetStorageByTempMatch().
const GetStorageByTempMatchBase = `SELECT * FROM storage WHERE temp_match="%s"`

//CreateStorageBase can be used with fmt.Sprintf() to get the Query for CreateStorage().
const CreateStorageBase = `INSERT INTO storage (personal_id, user_id, title, description, temp_match) ` +
	`VALUES(%d, %d, "%s", "%s", "%s")`

//UpdateStorageBase can be used with fmt.Sprintf() to get the Query for UpdateStorage().
const UpdateStorageBase = `UPDATE storage SET personal_id = %d, title = "%s", description = "%s", temp_match = "%s" WHERE id=%d`

//DeleteStorageBase can be used with fmt.Sprintf() to get the Query for DeleteStorage().
const DeleteStorageBase = `DELETE FROM storage WHERE user_id = %d AND personal_id=%d`

//GetStorageDishesBase can be used with fmt.Sprintf() to get the Query for GetStorageDishes().
const GetStorageDishesBase = `SELECT * FROM dish WHERE user_id = %d AND personal_id = %d`

//Repository interface is a contract for all the methods contained by this db.Repository object.
type Repository interface {
	GetDishes(int) (*dish.Dishes, fcerr.FCErr)
	GetDishByID(int, int) (*dish.Dish, fcerr.FCErr)
	GetDishByTempMatch(string) (*dish.Dish, fcerr.FCErr)
	GetPersonalDishCount(int) (int, fcerr.FCErr)
	CreateDish(dish.Dish) (*dish.Dish, fcerr.FCErr)
	UpdateDish(dish.Dish) fcerr.FCErr
	DeleteDish(int, int) fcerr.FCErr

	GetUsers() (*user.Users, fcerr.FCErr)
	GetUserByID(int) (*user.User, fcerr.FCErr)
	GetUserByEmail(string) (*user.User, fcerr.FCErr)
	GetUserByAlexa(string) (*user.User, fcerr.FCErr)
	GetUserByTempMatch(string) (*user.User, fcerr.FCErr)
	CreateUser(user.User) (*user.User, fcerr.FCErr)
	UpdateUser(user.User) (*user.User, fcerr.FCErr)
	DeleteUser(user.User) fcerr.FCErr

	GetStorages(int) (*storage.Storages, fcerr.FCErr)
	GetStorageByID(int, int) (*storage.Storage, fcerr.FCErr)
	GetStorageByTempMatch(string) (*storage.Storage, fcerr.FCErr)
	CreateStorage(storage.Storage) (*storage.Storage, fcerr.FCErr)
	UpdateStorage(storage.Storage) (*storage.Storage, fcerr.FCErr)
	DeleteStorage(int, int) fcerr.FCErr

	GetStorageDishes(int, int) (*dish.Dishes, fcerr.FCErr)
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

//NewRepositoryWithDB will get an instance of this type which satisfies the Repository interface.
func NewRepositoryWithDB(db *sql.DB) (Repository, fcerr.FCErr) {
	resultDB := repository{db}
	return &resultDB, nil
}

//GetDishes returns the list of all dishes in the database
func (repo *repository) GetDishes(userID int) (*dish.Dishes, fcerr.FCErr) {
	fmt.Println("now at the beginning of the db_repository GetDishes()")
	var resultDishes dish.Dishes
	getDishesQuery := fmt.Sprintf(GetDishesBase, userID)
	rows, err := repo.db.Query(getDishesQuery)
	fmt.Println("now after doing the Query:", getDishesQuery)
	if err != nil {
		fmt.Println("got an error on the Query:", err.Error())
		fcerr := fcerr.NewInternalServerError("Error while retrieving dishes from the database")
		return nil, fcerr
	}
	defer rows.Close()
	fmt.Println("now about to check the rows returned:")
	count := 0
	for rows.Next() {
		count++
		var currentDish dish.Dish
		fmt.Println("Inside the result set loop. currentDish:", currentDish)
		err := rows.Scan(&currentDish.DishID, &currentDish.PersonalDishID, &currentDish.UserID, &currentDish.StorageID, &currentDish.Title,
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

//GetDishByID (userID int, pID int) queries the mysql database for a dish the requesting user has with the given personal id.
func (repo *repository) GetDishByID(userID int, pID int) (*dish.Dish, fcerr.FCErr) {
	var resultingDish dish.Dish
	getDishByIDQuery := fmt.Sprintf(GetDishByIDBase, userID, pID)
	fmt.Println("about to run this query in GetDishByID:", getDishByIDQuery)

	rows, err := repo.db.Query(getDishByIDQuery)
	fmt.Println("now after doing the Query:", getDishByIDQuery)
	if err != nil {
		fmt.Println("got an error on the Query:", err.Error())
		fcerr := fcerr.NewInternalServerError("Error while retrieving dish from the database")
		return nil, fcerr
	}
	defer rows.Close()
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
		err := rows.Scan(&currentDish.DishID, &currentDish.PersonalDishID, &currentDish.UserID, &currentDish.StorageID, &currentDish.Title,
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
		err := rows.Scan(&currentDish.DishID, &currentDish.PersonalDishID, &currentDish.UserID, &currentDish.StorageID, &currentDish.Title,
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
	tMatch := generateTempMatch()
	createDishQuery := fmt.Sprintf(CreateDishBase, d.PersonalDishID, d.UserID, d.StorageID, d.Title, d.Description,
		d.CreatedDate, d.ExpireDate, d.Priority, d.DishType, d.Portions, tMatch)

	fmt.Println("About to run this Query on the database:\n", createDishQuery)

	_, err := repo.db.Query(createDishQuery)
	if err != nil {
		fmt.Println("got an error on the Query:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while inserting the dish into the database")
		return nil, fcerr
	}

	checkDish, err := repo.GetDishByTempMatch(tMatch)
	if err != nil {
		fmt.Println("Trying to CreateDish, seem to have hit a snag. Got an error when checking what we just put in: " + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while checking the dish that was created." +
			" Cannot verify if anything was entered to the Database")
		return nil, fcerr
	}

	return checkDish, nil
}

//UpdateDish takes a dish object and tries to update the existing dish in the database to match
func (repo *repository) UpdateDish(d dish.Dish) fcerr.FCErr {
	updateDishQuery := fmt.Sprintf(UpdateDishBase, d.PersonalDishID, d.StorageID, d.Title, d.Description,
		d.ExpireDate, d.Priority, d.DishType, d.Portions, d.DishID)

	fmt.Println("About to run this Query on the database:\n", updateDishQuery)

	_, err := repo.db.Query(updateDishQuery)
	if err != nil {
		fmt.Println("got an error on the query:" + err.Error())
		return fcerr.NewInternalServerError("Error while updating the dish in the database")
	}

	_, err2 := repo.GetDishByID(d.UserID, d.PersonalDishID)
	if err2 != nil {
		fmt.Println("got an error on the check query:")
		return fcerr.NewInternalServerError("Error while checking the dish that was created. Cannot verify if anything was updated in the Database")
	}

	return nil
}

//GetPersonalDishCount gets the number of dishes the given user has in the database
func (repo *repository) GetPersonalDishCount(userID int) (int, fcerr.FCErr) {
	getPersonalDishCountQuery := fmt.Sprintf(GetPersonalDishCountBase, userID)
	personalDishCountRow := repo.db.QueryRow(getPersonalDishCountQuery)
	var personalDishCount int
	err := personalDishCountRow.Scan(&personalDishCount)
	if err != nil {
		fmt.Println("got an error on the get personal count process:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while deleting the dish from the database - before delete command")
		return 0, fcerr
	}
	return personalDishCount, nil

}

//DeleteDish takes a requesting user and a personal dish id and tries to delete the dish
func (repo *repository) DeleteDish(userID int, pID int) fcerr.FCErr {
	personalDishCount, err := repo.GetPersonalDishCount(userID)
	if err != nil {
		return fcerr.NewInternalServerError("Error when Deleting the dish")
	}

	pDSstr := ""
	for i := pID + 1; i <= personalDishCount; i++ {
		if i == personalDishCount {
			pDSstr = pDSstr + strconv.Itoa(i)
		} else {
			pDSstr = pDSstr + strconv.Itoa(i) + ","
		}

	}

	deleteDishQuery := fmt.Sprintf(DeleteDishBase, userID, pID)

	_, err2 := repo.db.Query(deleteDishQuery)
	if err2 != nil {
		fmt.Println("got an error on the delete query:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while deleting the dish from the database")
		return fcerr
	}

	//TODO - fix hardcoding of user id 1 after we have the requesting user here
	decrementSomeDishesQuery := fmt.Sprintf(DecrementSomeDishesBase, 1, pDSstr)
	fmt.Println("about to run this query on the db:", decrementSomeDishesQuery)
	_, err3 := repo.db.Query(decrementSomeDishesQuery)
	if err3 != nil {
		fmt.Println("got an error while trying to decrement some dishes:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while deleting the dish from the database")
		return fcerr
	}

	returnedDish, err := repo.GetDishByID(userID, pID)
	if err == nil {
		fmt.Println("Expected an error here, but didn't get one!! Dish Title:" + returnedDish.Title)
		fcerr := fcerr.NewInternalServerError("Error while deleting the dish from the database, could not verify it was deleted.")
		return fcerr
	}

	return nil
}

//GetUsers queries the database and returns a slice of User objects
func (repo *repository) GetUsers() (*user.Users, fcerr.FCErr) {
	fmt.Println("now at the beginning of the db_repository GetUsers()")
	var resultingUsers user.Users
	getUsersQuery := fmt.Sprintf(GetUsersBase)
	rows, err := repo.db.Query(getUsersQuery)
	fmt.Println("now after doing the Query:", getUsersQuery)
	if err != nil {
		fmt.Println("got an error on the Query:", err.Error())
		fcerr := fcerr.NewInternalServerError("Error while retrieving users from the database")
		return nil, fcerr
	}
	defer rows.Close()
	fmt.Println("now about to check the rows returned:")
	count := 0
	for rows.Next() {
		count++
		var currentUser user.User
		fmt.Println("Inside the result set loop. currentDish:", currentUser)
		err := rows.Scan(&currentUser.UserID, &currentUser.Email, &currentUser.FirstName, &currentUser.LastName, &currentUser.FullName,
			&currentUser.CreatedDate, &currentUser.AccessToken, &currentUser.RefreshToken, &currentUser.AlexaUserID, &currentUser.Admin, &currentUser.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fmt.Println("&currentUser.UserID:", currentUser.UserID)
			fmt.Println("&currentUser.TempMatch:", currentUser.TempMatch)
			fcerr := fcerr.NewInternalServerError("Error while scanning the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current user scanned. currentUser:", currentUser)
		resultingUsers = append(resultingUsers, currentUser)

	}
	if count < 1 {
		fcerr := fcerr.NewNotFoundError("Database could not find any users")
		fmt.Println("Database could not find any users")
		return nil, fcerr
	}

	return &resultingUsers, nil
}

//GetUserByID gets a user from the database with the given ID.
func (repo *repository) GetUserByID(id int) (*user.User, fcerr.FCErr) {
	getUserByIDQuery := fmt.Sprintf(GetUserByIDBase, id)
	fmt.Println("About to run this Query on the database:\n", getUserByIDQuery)
	var resultingUser user.User

	rows, err := repo.db.Query(getUserByIDQuery)
	if err != nil {
		fmt.Println("got an error on the Query")
		fcerr := fcerr.NewInternalServerError("Error while retrieving user from the database")
		return nil, fcerr
	}
	defer rows.Close()
	fmt.Println("now about to check the rows returned:")
	count := 0
	for rows.Next() {
		count++
		if count > 1 {
			dberr := fcerr.NewInternalServerError("Database returned more than 1 row when only 1 was expected")
			return nil, dberr
		}
		var cUser user.User
		fmt.Println("Inside the result set loop. currentUser:", cUser)
		err := rows.Scan(&cUser.UserID, &cUser.Email, &cUser.FirstName, &cUser.LastName, &cUser.FullName,
			&cUser.CreatedDate, &cUser.AccessToken, &cUser.RefreshToken, &cUser.AlexaUserID, &cUser.Admin, &cUser.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fcerr := fcerr.NewInternalServerError("Error while scanning the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current user scanned. currentUser:", cUser)
		resultingUser = cUser

	}
	if count == 0 {
		fcerr := fcerr.NewNotFoundError("Database could not find a user with this ID")
		return nil, fcerr
	}
	return &resultingUser, nil
}

//GetUserByEmail gets a user from the database with the given Email.
func (repo *repository) GetUserByEmail(email string) (*user.User, fcerr.FCErr) {
	getUserByEmailQuery := fmt.Sprintf(GetUserByEmailBase, email)
	fmt.Println("About to run this Query on the database:\n", getUserByEmailQuery)
	var resultingUser user.User

	rows, err := repo.db.Query(getUserByEmailQuery)
	if err != nil {
		fmt.Println("got an error on the Query")
		fcerr := fcerr.NewInternalServerError("Error while retrieving user from the database")
		return nil, fcerr
	}
	defer rows.Close()
	fmt.Println("now about to check the rows returned:")
	count := 0
	for rows.Next() {
		count++
		if count > 1 {
			dberr := fcerr.NewInternalServerError("Database returned more than 1 row when only 1 was expected")
			return nil, dberr
		}
		var cUser user.User
		fmt.Println("Inside the result set loop. currentUser:", cUser)
		err := rows.Scan(&cUser.UserID, &cUser.Email, &cUser.FirstName, &cUser.LastName, &cUser.FullName,
			&cUser.CreatedDate, &cUser.AccessToken, &cUser.RefreshToken, &cUser.AlexaUserID, &cUser.Admin, &cUser.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fcerr := fcerr.NewInternalServerError("Error while scanning the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current user scanned. currentUser:", cUser)
		resultingUser = cUser

	}
	if count == 0 {
		fcerr := fcerr.NewNotFoundError("Database could not find a user with this Email")
		return nil, fcerr
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
		fcerr := fcerr.NewInternalServerError("Error while retrieving user from the database")
		return nil, fcerr
	}
	defer rows.Close()
	fmt.Println("now about to check the rows returned:")
	count := 0
	for rows.Next() {
		count++
		if count > 1 {
			dberr := fcerr.NewInternalServerError("Database returned more than 1 row when only 1 was expected")
			return nil, dberr
		}
		var cUser user.User
		fmt.Println("Inside the result set loop. currentUser:", cUser)
		err := rows.Scan(&cUser.UserID, &cUser.Email, &cUser.FirstName, &cUser.LastName, &cUser.FullName,
			&cUser.CreatedDate, &cUser.AccessToken, &cUser.RefreshToken, &cUser.AlexaUserID, &cUser.Admin, &cUser.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fcerr := fcerr.NewInternalServerError("Error while scanning the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current user scanned. currentUser:", cUser)
		resultingUser = cUser

	}
	if count == 0 {
		fcerr := fcerr.NewNotFoundError("Database could not find a user with this Alexa User ID")
		return nil, fcerr
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
		fcerr := fcerr.NewInternalServerError("Error while retrieving user from the database")
		return nil, fcerr
	}
	defer rows.Close()
	fmt.Println("now about to check the rows returned:")
	count := 0
	for rows.Next() {
		count++
		if count > 1 {
			dberr := fcerr.NewInternalServerError("Database returned more than 1 row when only 1 was expected")
			return nil, dberr
		}
		var cUser user.User
		fmt.Println("Inside the result set loop. currentUser:", cUser)
		err := rows.Scan(&cUser.UserID, &cUser.Email, &cUser.FirstName, &cUser.LastName, &cUser.FullName,
			&cUser.CreatedDate, &cUser.AccessToken, &cUser.RefreshToken, &cUser.AlexaUserID, &cUser.Admin, &cUser.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fcerr := fcerr.NewInternalServerError("Error while scanning the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current user scanned. currentUser:", cUser)
		resultingUser = cUser

	}
	if count == 0 {
		fcerr := fcerr.NewNotFoundError("Database could not find a user with this Temp Match")
		return nil, fcerr
	}
	return &resultingUser, nil
}

func (repo *repository) CreateUser(u user.User) (*user.User, fcerr.FCErr) {
	tMatch := generateTempMatch()
	createUserQuery := fmt.Sprintf(CreateUserBase, u.Email, u.FirstName, u.LastName, u.FullName,
		u.CreatedDate, u.AccessToken, u.RefreshToken, u.AlexaUserID, u.Admin, tMatch)

	fmt.Println("About to run this Query on the database:\n", createUserQuery)

	_, err := repo.db.Query(createUserQuery)
	if err != nil {
		fmt.Println("got an error on the Query:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while inserting the user into the database")
		return nil, fcerr
	}

	checkUser, err := repo.GetUserByTempMatch(tMatch)
	if err != nil {
		fmt.Println("Trying to CreateUser, seem to have hit a snag. Got an error when checking what we just put in: " + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while checking the user that was created." +
			" Cannot verify if anything was entered to the Database")
		return nil, fcerr
	}

	return checkUser, nil
}

//UpdateUser takes a user object and tries to update the existing user in the database to match
func (repo *repository) UpdateUser(u user.User) (*user.User, fcerr.FCErr) {
	updateUserQuery := fmt.Sprintf(UpdateUserBase, u.Email, u.FirstName, u.LastName,
		u.FullName, u.AccessToken, u.RefreshToken, u.AlexaUserID, u.TempMatch, u.UserID)

	fmt.Println("About to run this Query on the database:\n", updateUserQuery)

	_, err := repo.db.Query(updateUserQuery)
	if err != nil {
		fmt.Println("got an error on the query:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while updating the user in the database")
		return nil, fcerr
	}

	checkDish, err := repo.GetUserByID(u.UserID)
	if err != nil {
		fmt.Println("got an error on the check query:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while checking the user that was created." +
			" Cannot verify if anything was updated in the Database")
		return nil, fcerr
	}

	return checkDish, nil
}

//DeleteUser takes a user object and tries to delete the existing user from the database
func (repo *repository) DeleteUser(u user.User) fcerr.FCErr {
	deleteUserQuery := fmt.Sprintf(DeleteUserBase, u.UserID)

	_, err := repo.db.Query(deleteUserQuery)
	if err != nil {
		fmt.Println("got an error on the delete query:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while deleting the user from the database")
		return fcerr

	}

	returnedUser, err := repo.GetUserByID(u.UserID)
	if err == nil {
		fmt.Println("Expected an error here, but didn't get one!! User Email:" + returnedUser.Email)
		fcerr := fcerr.NewInternalServerError("Error while deleting the user from the database, could not verify it was deleted.")
		return fcerr
	}

	return nil
}

//GetStorage takes an int of a user id and returns the list of storage objects owned by that user.
func (repo *repository) GetStorages(userID int) (*storage.Storages, fcerr.FCErr) {
	fmt.Println("now at the beginning of the db_repository GetStoragesByUser()")
	var resultingStorages storage.Storages
	getStoragesQuery := fmt.Sprintf(GetStoragesBase, userID)
	rows, err := repo.db.Query(getStoragesQuery)
	fmt.Println("now after doing the Query:", getStoragesQuery)
	if err != nil {
		fmt.Println("got an error on the Query:", err.Error())
		fcerr := fcerr.NewInternalServerError("Error while retrieving storage units from the database")
		return nil, fcerr
	}
	defer rows.Close()
	fmt.Println("now about to check the rows returned:")
	count := 0
	for rows.Next() {
		count++
		var currentStorage storage.Storage
		fmt.Println("Inside the result set loop. currentStorage:", currentStorage)
		err := rows.Scan(&currentStorage.StorageID, &currentStorage.UserID, &currentStorage.Title, &currentStorage.Description, &currentStorage.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fmt.Println("&currentStorage.StorageID:", currentStorage.StorageID)
			fmt.Println("&currentStorage.TempMatch:", currentStorage.TempMatch)
			fcerr := fcerr.NewInternalServerError("Error while scanning the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current storage scanned. currentStorage:", currentStorage)
		resultingStorages = append(resultingStorages, currentStorage)

	}
	if count < 1 {
		fmt.Println("Database could not find any storage units for this user")
		fcerr := fcerr.NewNotFoundError("Database could not find any storage units for this user")
		return nil, fcerr
	}

	return &resultingStorages, nil
}

//GetStorageByID takes an int and queries the mysql database for a storage with this id.
func (repo *repository) GetStorageByID(userID int, pID int) (*storage.Storage, fcerr.FCErr) {
	getStorageByIDQuery := fmt.Sprintf(GetStorageByIDBase, userID, pID)
	fmt.Println("About to run this Query on the database:\n", getStorageByIDQuery)
	var resultingStorage storage.Storage

	rows, err := repo.db.Query(getStorageByIDQuery)
	if err != nil {
		fmt.Println("got an error on the Query")
		fcerr := fcerr.NewInternalServerError("Error while retrieving storage unit from the database")
		return nil, fcerr
	}
	defer rows.Close()
	fmt.Println("now about to check the rows returned:")
	count := 0
	for rows.Next() {
		count++
		if count > 1 {
			dberr := fcerr.NewInternalServerError("Database returned more than 1 row when only 1 was expected")
			return nil, dberr
		}
		var cStorage storage.Storage
		fmt.Println("Inside the result set loop. currentStorage:", cStorage)
		err := rows.Scan(&cStorage.StorageID, &cStorage.UserID, &cStorage.Title, &cStorage.Description, &cStorage.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fcerr := fcerr.NewInternalServerError("Error while scanning the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current storage unit scanned. currentStorage:", cStorage)
		resultingStorage = cStorage

	}
	if count == 0 {
		fcerr := fcerr.NewNotFoundError("Database could not find a storage unit with this ID")
		return nil, fcerr
	}
	return &resultingStorage, nil
}

//GetStorageByTempMatch takes a string and queries the mysql database for a storage with this temp_match.
func (repo *repository) GetStorageByTempMatch(tM string) (*storage.Storage, fcerr.FCErr) {
	getStorageByIDQuery := fmt.Sprintf(GetStorageByTempMatchBase, tM)
	fmt.Println("About to run this Query on the database:\n", getStorageByIDQuery)
	var resultingStorage storage.Storage

	rows, err := repo.db.Query(getStorageByIDQuery)
	if err != nil {
		fmt.Println("got an error on the Query")
		fcerr := fcerr.NewInternalServerError("Error while retrieving storage unit from the database")
		return nil, fcerr
	}
	defer rows.Close()
	fmt.Println("now about to check the rows returned:")
	count := 0
	for rows.Next() {
		count++
		if count > 1 {
			dberr := fcerr.NewInternalServerError("Database returned more than 1 row when only 1 was expected")
			return nil, dberr
		}
		var cStorage storage.Storage
		fmt.Println("Inside the result set loop. currentStorage:", cStorage)
		err := rows.Scan(&cStorage.StorageID, &cStorage.UserID, &cStorage.Title, &cStorage.Description, &cStorage.TempMatch)
		if err != nil {
			fmt.Println("got an error from the rows.Scan.")
			fcerr := fcerr.NewInternalServerError("Error while scanning the result from the database")
			return nil, fcerr
		}
		fmt.Println("now after the current storage unit scanned. currentStorage:", cStorage)
		resultingStorage = cStorage

	}
	if count == 0 {
		fcerr := fcerr.NewNotFoundError("Database could not find a storage unit with this ID")
		return nil, fcerr
	}
	return &resultingStorage, nil
}

//CreateStorage takes a storage object and tries to add it to the database
func (repo *repository) CreateStorage(s storage.Storage) (*storage.Storage, fcerr.FCErr) {
	tMatch := generateTempMatch()
	createStorageQuery := fmt.Sprintf(CreateStorageBase, s.UserID, s.Title, s.Description, tMatch)

	fmt.Println("About to run this Query on the database:\n", createStorageQuery)

	_, err := repo.db.Query(createStorageQuery)
	if err != nil {
		fmt.Println("got an error on the Query:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while inserting the storage unit into the database")
		return nil, fcerr
	}

	checkStorage, err := repo.GetStorageByTempMatch(tMatch)
	if err != nil {
		fmt.Println("Trying to CreateStorage, seem to have hit a snag. Got an error when checking what we just put in: " + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while checking the storage unit that was created." +
			" Cannot verify if anything was entered to the Database")
		return nil, fcerr
	}

	return checkStorage, nil
}

//UpdateStorage takes a storage object and tries to update the existing storage in the database to match
func (repo *repository) UpdateStorage(s storage.Storage) (*storage.Storage, fcerr.FCErr) {
	updateStorageQuery := fmt.Sprintf(UpdateStorageBase, s.Title, s.Description, s.TempMatch, s.StorageID)

	fmt.Println("About to run this Query on the database:\n", updateStorageQuery)

	_, err := repo.db.Query(updateStorageQuery)
	if err != nil {
		fmt.Println("got an error on the query:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while updating the storage unit in the database")
		return nil, fcerr
	}

	checkStorage, err := repo.GetStorageByID(s.UserID, s.PersonalID)
	if err != nil {
		fmt.Println("got an error on the check query:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while checking the storage unit that was created." +
			" Cannot verify if anything was updated in the Database")
		return nil, fcerr
	}

	return checkStorage, nil
}

//DeleteStorage takes a storage object and tries to delete the existing storage from the database
func (repo *repository) DeleteStorage(userID int, pID int) fcerr.FCErr {
	deleteStorageQuery := fmt.Sprintf(DeleteStorageBase, userID, pID)

	_, err := repo.db.Query(deleteStorageQuery)
	if err != nil {
		fmt.Println("got an error on the delete query:" + err.Error())
		fcerr := fcerr.NewInternalServerError("Error while deleting the storage unit from the database")
		return fcerr

	}

	returnedStorage, err := repo.GetStorageByID(userID, pID)
	if err == nil {
		fmt.Println("Expected an error here, but didn't get one!! Storage ID:", returnedStorage.StorageID)
		fcerr := fcerr.NewInternalServerError("Error while deleting the storage unit from the database, could not verify it was deleted.")
		return fcerr
	}

	return nil
}

//GetStorageDishes takes a storage object and tries to update the existing storage in the database to match
func (repo *repository) GetStorageDishes(userID int, storagePID int) (*dish.Dishes, fcerr.FCErr) {
	fmt.Println("now at the beginning of the db_repository GetStorageDishes()")
	var resultDishes dish.Dishes
	getStorageDishesQuery := fmt.Sprintf(GetStorageDishesBase, userID, storagePID)
	rows, err := repo.db.Query(getStorageDishesQuery)
	fmt.Println("now after doing the Query:", getStorageDishesQuery)
	if err != nil {
		fmt.Println("got an error on the Query:", err.Error())
		fcerr := fcerr.NewInternalServerError("Error while retrieving dishes from the database")
		return nil, fcerr
	}
	defer rows.Close()
	fmt.Println("now about to check the rows returned:")
	count := 0
	for rows.Next() {
		count++
		var currentDish dish.Dish
		fmt.Println("Inside the result set loop. currentDish:", currentDish)
		err := rows.Scan(&currentDish.DishID, &currentDish.PersonalDishID, &currentDish.UserID, &currentDish.StorageID, &currentDish.Title,
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
		fcerr := fcerr.NewNotFoundError("Database could not find any dishes that belong to this storage unit")
		fmt.Println("Database could not find any storage dishes")
		return nil, fcerr
	}

	return &resultDishes, nil
}

func generateTempMatch() string {
	n := make([]byte, 15)
	rand.Read(n)
	fmt.Println("New way:", base64.URLEncoding.EncodeToString(n))

	return base64.URLEncoding.EncodeToString(n)

}
