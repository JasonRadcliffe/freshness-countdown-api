package db

import (
	"errors"
	"net/http"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDb_NewRepository_CantConnect(t *testing.T) {
	_, err := NewRepository("")

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	assert.Equal(t, "Error while connecting to the mysql database", err.Message())

}

func TestDb_GetDishes(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(1, 1, 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "").
		AddRow(1, 2, 3, "Peas", "Some peas we got at the store", "2007-01-02T15:04:05", "2021-10-13T08:00", 1, "", -1, "")

	mock.ExpectQuery("SELECT * FROM dish").WillReturnRows(rows)

	resultingDishes, err := repo.GetDishes()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(*resultingDishes))

	resultingDish1 := (*resultingDishes)[0]
	resultingDish2 := (*resultingDishes)[1]

	assert.Equal(t, "Carrots", resultingDish1.Title)
	assert.Equal(t, "Peas", resultingDish2.Title)

	assert.Equal(t, 1, resultingDish1.UserID)
	assert.Equal(t, 2, resultingDish2.UserID)
}

func TestDb_GetDishes_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"})

	mock.ExpectQuery("SELECT * FROM dish").WillReturnRows(rows)

	resultingDishes, err := repo.GetDishes()

	assert.NotNil(t, err)
	assert.Nil(t, resultingDishes)
	assert.Equal(t, "Database could not find any dishes", err.Message())
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestDb_GetDishes_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	newErr := errors.New("Error 1146: Table 'food_db.dishs' doesn't exist")
	mock.ExpectQuery("SELECT * FROM dishs").WillReturnError(newErr)
	resultingDishes, err := repo.GetDishes()

	assert.Nil(t, resultingDishes)
	assert.NotNil(t, err)
	assert.Equal(t, "Error while retrieving dishes from the database", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetDishes_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow("SHOULDBEINT", 1, 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "")

	mock.ExpectQuery("SELECT * FROM dish").WillReturnRows(rows)

	resultingDishes, err := repo.GetDishes()

	assert.Nil(t, resultingDishes)
	assert.NotNil(t, err)
	assert.Equal(t, "Error while scanning the result from the database", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetDishByID(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(1, 2, 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "")

	mock.ExpectQuery("SELECT * FROM dish WHERE id = 1").WillReturnRows(rows)

	resultingDish, err := repo.GetDishByID(1)

	assert.Nil(t, err)

	assert.Equal(t, 1, resultingDish.DishID)
	assert.Equal(t, "Carrots", resultingDish.Title)
	assert.Equal(t, 2, resultingDish.UserID)
	assert.Equal(t, 3, resultingDish.StorageID)
}

func TestDb_GetDishByID_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(1, 2, 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "")

	mock.ExpectQuery("SELECT * FROM dish WHERE id = 1").WillReturnRows(rows)

	resultingDish, err := repo.GetDishByID(1)

	assert.Nil(t, err)

	assert.Equal(t, 1, resultingDish.DishID)
	assert.Equal(t, "Carrots", resultingDish.Title)
	assert.Equal(t, 2, resultingDish.UserID)
	assert.Equal(t, 3, resultingDish.StorageID)
}

func TestDb_GetDishByID_FoundMultiple(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_GetDishByTempMatch(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_CreateDish(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_UpdateDish(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_DeleteDish(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_GetUsers(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_GetUserByID(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_GetUserByEmail(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_GetUserByAlexa(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_GetUserByTempMatch(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_CreateUser(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_UpdateUser(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_DeleteUser(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_GetStorage(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_GetStorageByID(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_CreateStorage(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_UpdateStorage(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_DeleteStorage(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_GetStorageDishes(t *testing.T) {
	assert.Equal(t, "", "")
}
