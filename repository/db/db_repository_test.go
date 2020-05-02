package db

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/jasonradcliffe/freshness-countdown-api/domain/dish"

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
		"expire_date", "priority", "dish_type", "portions", "temp_match"})

	mock.ExpectQuery("SELECT * FROM dish WHERE id = 1").WillReturnRows(rows)

	resultingDish, err := repo.GetDishByID(1)

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	assert.Equal(t, http.StatusNotFound, err.Status())
	assert.Equal(t, "Database could not find a dish with this ID", err.Message())
}

func TestDb_GetDishByID_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(1, "SHOULD BE INT", 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "")

	mock.ExpectQuery("SELECT * FROM dish WHERE id = 1").WillReturnRows(rows)

	resultingDish, err := repo.GetDishByID(1)

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	assert.Equal(t, "Error while scanning the result from the database", err.Message())
}

func TestDb_GetDishByID_FoundMultiple(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(1, 2, 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "").
		AddRow(1, 2, 3, "Carrots", "Some carrots we got at the store a second time", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "")

	mock.ExpectQuery("SELECT * FROM dish WHERE id = 1").WillReturnRows(rows)

	resultingDish, err := repo.GetDishByID(1)

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	assert.Equal(t, "Database returned more than 1 row when only 1 was expected", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetDishByTempMatch(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(1, 2, 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "9r842da351")

	mock.ExpectQuery(fmt.Sprintf(GetDishByTempMatchBase, "9r842da351")).WillReturnRows(rows)

	resultingDish, err := repo.GetDishByTempMatch("9r842da351")

	assert.Nil(t, err)

	assert.Equal(t, 1, resultingDish.DishID)
	assert.Equal(t, "Carrots", resultingDish.Title)
	assert.Equal(t, 2, resultingDish.UserID)
	assert.Equal(t, 3, resultingDish.StorageID)
	assert.Equal(t, "9r842da351", resultingDish.TempMatch)
}

func TestDb_GetDishByTempMatch_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"})

	mock.ExpectQuery(`Select * FROM dish WHERE temp_match = "9r842da351"`).WillReturnRows(rows)

	resultingDish, err := repo.GetDishByTempMatch("9r842da351")

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	assert.Equal(t, http.StatusNotFound, err.Status())
	assert.Equal(t, "Database could not find a dish with this temp match", err.Message())

}

func TestDb_GetDishByTempMatch_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(1, "SHOULD BE INT", 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "")

	mock.ExpectQuery(`Select * FROM dish WHERE temp_match = "9r842da351"`).WillReturnRows(rows)

	resultingDish, err := repo.GetDishByTempMatch("9r842da351")

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	assert.Equal(t, "Error while scanning the result from the database", err.Message())

}

func TestDb_GetDishByTempMatch_FoundMultiple(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(1, 2, 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "9r842da351").
		AddRow(4, 2, 3, "Carrots", "Some carrots we got at the store a second time", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "9r842da351")

	mock.ExpectQuery(`Select * FROM dish WHERE temp_match = "9r842da351"`).WillReturnRows(rows)

	resultingDish, err := repo.GetDishByTempMatch("9r842da351")

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	assert.Equal(t, "Database returned more than 1 row when only 1 was expected", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_CreateDish(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nD := &dish.Dish{
		UserID:      2,
		StorageID:   3,
		Title:       "Carrots",
		Description: "Some carrots we got at the store",
		CreatedDate: "2006-01-02T15:04:05",
		ExpireDate:  "2020-10-13T08:00",
		Priority:    "",
		DishType:    "",
		Portions:    -1,
		TempMatch:   "9r842d3a351",
	}

	createRows := sqlmock.NewRows([]string{""})

	getRows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(5, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(CreateDishBase, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
		nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)).
		WillReturnRows(createRows)

	mock.ExpectQuery(fmt.Sprintf(GetDishByTempMatchBase, nD.TempMatch)).WillReturnRows(getRows)

	returnedDish, err := repo.CreateDish(*nD)

	assert.Nil(t, err)

	assert.NotNil(t, returnedDish)
	assert.Equal(t, nD.Title, returnedDish.Title)
}

func TestDb_CreateDish_InsertError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nD := &dish.Dish{
		UserID:      2,
		StorageID:   3,
		Title:       "Carrots",
		Description: "Some carrots we got at the store",
		CreatedDate: "2006-01-02T15:04:05",
		ExpireDate:  "2020-10-13T08:00",
		Priority:    "",
		DishType:    "",
		Portions:    -1,
		TempMatch:   "9r842d3a351",
	}

	mock.ExpectQuery(fmt.Sprintf(CreateDishBase, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)).
		WillReturnError(errors.New("not possible"))

	returnedDish, err := repo.CreateDish(*nD)

	assert.NotNil(t, err)
	assert.Nil(t, returnedDish)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	assert.Equal(t, "Error while inserting the dish into the database", err.Message())

	assert.NotNil(t, returnedDish)
	assert.Equal(t, nD.Title, returnedDish.Title)
}

func TestDb_CreateDish_CheckError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nD := &dish.Dish{
		UserID:      2,
		StorageID:   3,
		Title:       "Carrots",
		Description: "Some carrots we got at the store",
		CreatedDate: "2006-01-02T15:04:05",
		ExpireDate:  "2020-10-13T08:00",
		Priority:    "",
		DishType:    "",
		Portions:    -1,
		TempMatch:   "9r842d3a351",
	}

	createRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(fmt.Sprintf(CreateDishBase, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
		nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)).
		WillReturnRows(createRows)

	mock.ExpectQuery(fmt.Sprintf(GetDishByTempMatchBase, nD.TempMatch)).
		WillReturnError(errors.New("not possible"))

	returnedDish, err := repo.CreateDish(*nD)

	assert.NotNil(t, err)
	assert.Nil(t, returnedDish)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	assert.Equal(t, "Error while checking the dish that was created."+
		" Cannot verify if anything was entered to the Database", err.Message())
}

func TestDb_UpdateDish(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nD := &dish.Dish{
		DishID:      2,
		UserID:      2,
		StorageID:   3,
		Title:       "Carrots",
		Description: "Some carrots we got at the store",
		CreatedDate: "2006-01-02T15:04:05",
		ExpireDate:  "2020-10-13T08:00",
		Priority:    "",
		DishType:    "",
		Portions:    -1,
		TempMatch:   "9r842d3a351",
	}

	createRows := sqlmock.NewRows([]string{""})

	getRows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(2, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(UpdateDishBase, nD.StorageID, nD.Title,
		nD.Description, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.DishID)).
		WillReturnRows(createRows)

	mock.ExpectQuery(fmt.Sprintf(GetDishByIDBase, nD.DishID)).WillReturnRows(getRows)

	returnedDish, err := repo.UpdateDish(*nD)

	assert.Nil(t, err)
	assert.NotNil(t, returnedDish)

	assert.Equal(t, nD.DishID, returnedDish.DishID)
	assert.Equal(t, nD.StorageID, returnedDish.StorageID)
	assert.Equal(t, nD.Title, returnedDish.Title)
	assert.Equal(t, nD.Description, returnedDish.Description)
	assert.Equal(t, nD.ExpireDate, returnedDish.ExpireDate)
	assert.Equal(t, nD.Priority, returnedDish.Priority)
	assert.Equal(t, nD.DishType, returnedDish.DishType)
	assert.Equal(t, nD.Portions, returnedDish.Portions)

}

func TestDb_UpdateDish_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nD := &dish.Dish{
		DishID:      2,
		UserID:      2,
		StorageID:   3,
		Title:       "Carrots",
		Description: "Some carrots we got at the store",
		CreatedDate: "2006-01-02T15:04:05",
		ExpireDate:  "2020-10-13T08:00",
		Priority:    "",
		DishType:    "",
		Portions:    -1,
		TempMatch:   "9r842d3a351",
	}

	mock.ExpectQuery(fmt.Sprintf(UpdateDishBase, nD.StorageID, nD.Title,
		nD.Description, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.DishID)).
		WillReturnError(errors.New("database error"))

	returnedDish, err := repo.UpdateDish(*nD)

	assert.Nil(t, returnedDish)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	assert.Equal(t, "Error while updating the dish in the database", err.Message())
}

func TestDb_UpdateDish_CheckError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nD := &dish.Dish{
		DishID:      2,
		UserID:      2,
		StorageID:   3,
		Title:       "Carrots",
		Description: "Some carrots we got at the store",
		CreatedDate: "2006-01-02T15:04:05",
		ExpireDate:  "2020-10-13T08:00",
		Priority:    "",
		DishType:    "",
		Portions:    -1,
		TempMatch:   "9r842d3a351",
	}

	createRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(fmt.Sprintf(UpdateDishBase, nD.StorageID, nD.Title,
		nD.Description, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.DishID)).
		WillReturnRows(createRows)

	mock.ExpectQuery(fmt.Sprintf(GetDishByIDBase, nD.DishID)).WillReturnError(errors.New("database error"))

	returnedDish, err := repo.UpdateDish(*nD)

	assert.Nil(t, returnedDish)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	assert.Equal(t, "Error while checking the dish that was created."+
		" Cannot verify if anything was updated in the Database", err.Message())

}

func TestDb_DeleteDish(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nD := &dish.Dish{
		DishID:      2,
		UserID:      2,
		StorageID:   3,
		Title:       "Carrots",
		Description: "Some carrots we got at the store",
		CreatedDate: "2006-01-02T15:04:05",
		ExpireDate:  "2020-10-13T08:00",
		Priority:    "",
		DishType:    "",
		Portions:    -1,
		TempMatch:   "9r842d3a351",
	}

	deleteRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(fmt.Sprintf(DeleteDishBase, nD.DishID)).WillReturnRows(deleteRows)

	err := repo.DeleteDish(*nD)

	assert.Nil(t, err)
}

func TestDb_DeleteDish_QueryError(t *testing.T) {
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
