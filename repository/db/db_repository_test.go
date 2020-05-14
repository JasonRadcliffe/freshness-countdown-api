package db

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/domain/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/domain/user"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var nD = &dish.Dish{
	DishID:         360,
	PersonalDishID: 2,
	UserID:         2,
	StorageID:      3,
	Title:          "Carrots",
	Description:    "Some carrots we got at the store",
	CreatedDate:    "2006-01-02T15:04:05",
	ExpireDate:     "2020-10-13T08:00",
	Priority:       "",
	DishType:       "",
	Portions:       -1,
	TempMatch:      "9r842d3a351",
}

var nDex = &dish.Dish{
	DishID:         4,
	PersonalDishID: 1,
	UserID:         2,
	StorageID:      3,
	Title:          "Old Carrots",
	Description:    "Some carrots we got at the store last year",
	CreatedDate:    "2006-01-02T15:04:05",
	ExpireDate:     "2019-10-13T08:00",
	Priority:       "",
	DishType:       "",
	Portions:       -1,
	TempMatch:      "9r842d3a351",
}

var nU = &user.User{
	UserID:       2,
	Email:        "nothing@gmail.com",
	FirstName:    "Bob",
	LastName:     "Nothing",
	FullName:     "Bob Nothing",
	CreatedDate:  "2016-01-02T15:04:05",
	AccessToken:  "ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
	RefreshToken: "105i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
	AlexaUserID:  "qwertyuiop",
	Admin:        false,
	TempMatch:    "1v842d234523a",
}

var nS = &storage.Storage{
	StorageID:   5,
	PersonalID:  1,
	UserID:      2,
	Title:       "Fridge",
	Description: "The main fridge in the house",
	TempMatch:   "Eb2iev8zpxgy-dxe",
}

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

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description,
			nD.CreatedDate, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch).
		AddRow(nD.DishID+200, nD.PersonalDishID+1, nD.UserID, nD.StorageID, nD.Title, nD.Description,
			nD.CreatedDate, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetDishesBase, nU.UserID)).WillReturnRows(rows)

	resultingDishes, err := repo.GetDishes(nU.UserID)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(*resultingDishes))

	resultingDish1 := (*resultingDishes)[0]

	assert.Equal(t, nD, &resultingDish1)
}

func TestDb_GetDishes_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(GetDishesBase, nU.UserID)).WillReturnRows(rows)

	resultingDishes, err := repo.GetDishes(nU.UserID)

	assert.NotNil(t, err)
	assert.Nil(t, resultingDishes)
	//assert.Equal(t, "Database could not find any dishes", err.Message())
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestDb_GetDishes_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(GetDishesBase, nU.UserID)).WillReturnError(errors.New("database error"))
	resultingDishes, err := repo.GetDishes(nU.UserID)

	assert.Nil(t, resultingDishes)
	assert.NotNil(t, err)
	//assert.Equal(t, "Error while retrieving dishes from the database", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetDishes_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow("SHOULDBEINT", 1, 1, 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "")

	mock.ExpectQuery(fmt.Sprintf(GetDishesBase, nU.UserID)).WillReturnRows(rows)

	resultingDishes, err := repo.GetDishes(nU.UserID)

	assert.Nil(t, resultingDishes)
	assert.NotNil(t, err)
	//assert.Equal(t, "Error while scanning the result from the database", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetDishByID(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description,
			nD.CreatedDate, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetDishByIDBase, nD.UserID, nD.PersonalDishID)).WillReturnRows(rows)

	resultingDish, err := repo.GetDishByID(nD.UserID, nD.PersonalDishID)

	assert.Nil(t, err)

	assert.Equal(t, nD, resultingDish)
}

func TestDb_GetDishByID_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(GetDishByIDBase, nD.UserID, nD.PersonalDishID)).WillReturnError(errors.New("database error"))

	resultingDish, err := repo.GetDishByID(nD.UserID, nD.PersonalDishID)

	assert.Nil(t, resultingDish)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while retrieving dish from the database", err.Message())
}

func TestDb_GetDishByID_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(GetDishByIDBase, nD.UserID, nD.PersonalDishID)).WillReturnRows(rows)

	resultingDish, err := repo.GetDishByID(nD.UserID, nD.PersonalDishID)

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	assert.Equal(t, http.StatusNotFound, err.Status())
	//assert.Equal(t, "Database could not find a dish with this ID", err.Message())
}

func TestDb_GetDishByID_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, "SHOULDBEINT", nD.Title, nD.Description,
			nD.CreatedDate, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetDishByIDBase, nD.UserID, nD.PersonalDishID)).WillReturnRows(rows)

	resultingDish, err := repo.GetDishByID(nD.UserID, nD.PersonalDishID)

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while scanning the result from the database", err.Message())
}

func TestDb_GetDishByID_FoundMultiple(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description,
			nD.CreatedDate, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch).
		AddRow(nD.DishID+1, nD.PersonalDishID+1, nD.UserID, nD.StorageID, nD.Title, nD.Description,
			nD.CreatedDate, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetDishByIDBase, nD.UserID, nD.PersonalDishID)).WillReturnRows(rows)

	resultingDish, err := repo.GetDishByID(nD.UserID, nD.PersonalDishID)

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	//assert.Equal(t, "Database returned more than 1 row when only 1 was expected", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetDishByTempMatch(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(1, 1, 2, 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "9r842da351")

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

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"})

	mock.ExpectQuery(`SELECT * FROM dish WHERE temp_match = "9r842da351"`).WillReturnRows(rows)

	resultingDish, err := repo.GetDishByTempMatch("9r842da351")

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	assert.Equal(t, http.StatusNotFound, err.Status())

}

func TestDb_GetDishByTempMatch_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(1, 2, "SHOULD BE INT", 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "")

	mock.ExpectQuery(`SELECT * FROM dish WHERE temp_match = "9r842da351"`).WillReturnRows(rows)

	resultingDish, err := repo.GetDishByTempMatch("9r842da351")

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while scanning the result from the database", err.Message())

}

func TestDb_GetDishByTempMatch_FoundMultiple(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(1, 1, 2, 3, "Carrots", "Some carrots we got at the store", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "9r842da351").
		AddRow(4, 1, 2, 3, "Carrots", "Some carrots we got at the store a second time", "2006-01-02T15:04:05", "2020-10-13T08:00", 1, "", -1, "9r842da351")

	mock.ExpectQuery(`SELECT * FROM dish WHERE temp_match = "9r842da351"`).WillReturnRows(rows)

	resultingDish, err := repo.GetDishByTempMatch("9r842da351")

	assert.NotNil(t, err)
	assert.Nil(t, resultingDish)

	//assert.Equal(t, "Database returned more than 1 row when only 1 was expected", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_CreateDish(t *testing.T) {

	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nD := &dish.Dish{
		PersonalDishID: 1,
		UserID:         2,
		StorageID:      3,
		Title:          "Carrots",
		Description:    "Some carrots we got at the store",
		CreatedDate:    "2006-01-02T15:04:05",
		ExpireDate:     "2020-10-13T08:00",
		Priority:       "",
		DishType:       "",
		Portions:       -1,
		TempMatch:      "9r842d3a351",
	}

	createRows := sqlmock.NewRows([]string{""})

	getRows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(5, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`INSERT INTO dish \(personal_id, user_id, storage_id, title, description, created_date, expire_date, priority, dish_type, portions, temp_match\) VALUES\(1, 2, 3, ".+", ".+", ".+", ".+", "", "", -1, ".+"\)`).
		WillReturnRows(createRows)

	mock.ExpectQuery(`SELECT \* FROM dish WHERE temp_match = ".+"`).
		WillReturnRows(getRows)

	returnedDish, err := repo.CreateDish(*nD)

	assert.Nil(t, err)

	assert.NotNil(t, returnedDish)

	assert.Equal(t, nD.Title, returnedDish.Title)
}

func TestDb_CreateDish_InsertError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nD := &dish.Dish{
		PersonalDishID: 1,
		UserID:         2,
		StorageID:      3,
		Title:          "Carrots",
		Description:    "Some carrots we got at the store",
		CreatedDate:    "2006-01-02T15:04:05",
		ExpireDate:     "2020-10-13T08:00",
		Priority:       "",
		DishType:       "",
		Portions:       -1,
		TempMatch:      "9r842d3a351",
	}

	mock.ExpectQuery(`INSERT INTO dish \(personal_id, user_id, storage_id, title, description, created_date, expire_date, priority, dish_type, portions, temp_match\) VALUES\(2, 3, ".+", ".+", ".+", ".+", "", "", -1, ".+"\)`).
		WillReturnError(errors.New("not possible"))

	returnedDish, err := repo.CreateDish(*nD)

	assert.NotNil(t, err)
	assert.Nil(t, returnedDish)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while inserting the dish into the database", err.Message())
}

func TestDb_CreateDish_CheckError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nD := &dish.Dish{
		PersonalDishID: 1,
		UserID:         2,
		StorageID:      3,
		Title:          "Carrots",
		Description:    "Some carrots we got at the store",
		CreatedDate:    "2006-01-02T15:04:05",
		ExpireDate:     "2020-10-13T08:00",
		Priority:       "",
		DishType:       "",
		Portions:       -1,
		TempMatch:      "9r842d3a351",
	}

	createRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(`INSERT INTO dish \(personal_id, user_id, storage_id, title, description, created_date, expire_date, priority, dish_type, portions, temp_match\) VALUES\(2, 3, ".+", ".+", ".+", ".+", "", "", -1, ".+"\)`).
		WillReturnRows(createRows)

	mock.ExpectQuery(fmt.Sprintf(GetDishByTempMatchBase, `.+`)).
		WillReturnError(errors.New("not possible"))

	returnedDish, err := repo.CreateDish(*nD)

	assert.NotNil(t, err)
	assert.Nil(t, returnedDish)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_UpdateDish(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	createRows := sqlmock.NewRows([]string{""})

	getRows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(2, 1, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(UpdateDishBase, nD.PersonalDishID, nD.StorageID, nD.Title,
		nD.Description, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.DishID)).
		WillReturnRows(createRows)

	mock.ExpectQuery(fmt.Sprintf(GetDishByIDBase, nD.UserID, nD.PersonalDishID)).WillReturnRows(getRows)

	err := repo.UpdateDish(*nD)

	assert.Nil(t, err)

}

func TestDb_UpdateDish_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(UpdateDishBase, nD.PersonalDishID, nD.StorageID, nD.Title,
		nD.Description, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.DishID)).
		WillReturnError(errors.New("database error"))

	err := repo.UpdateDish(*nD)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_UpdateDish_CheckError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	createRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(fmt.Sprintf(UpdateDishBase, nD.PersonalDishID, nD.StorageID, nD.Title,
		nD.Description, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.DishID)).
		WillReturnRows(createRows)

	mock.ExpectQuery(fmt.Sprintf(GetDishByIDBase, nD.UserID, nD.PersonalDishID)).WillReturnError(errors.New("database error"))

	err := repo.UpdateDish(*nD)

	assert.NotNil(t, err)
	//assert.Equal(t, http.StatusInternalServerError, err.Status())
	assert.Equal(t, "", "")
}

func TestDb_DeleteDish(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	countRow := sqlmock.NewRows([]string{"COUNT(*)"}).
		AddRow(3)

	deleteRows := sqlmock.NewRows([]string{""})

	updateRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(fmt.Sprintf(GetPersonalDishCountBase, nD.UserID)).WillReturnRows(countRow)

	mock.ExpectQuery(fmt.Sprintf(DeleteDishBase, nD.UserID, nD.PersonalDishID)).WillReturnRows(deleteRows)

	mock.ExpectQuery(`UPDATE dish SET personal_id = personal_id - 1 WHERE user_id = 1 AND personal_id IN(3)`).WillReturnRows(updateRows)

	err := repo.DeleteDish(nU.UserID, nD.PersonalDishID)

	assert.Nil(t, err)
}

func TestDb_DeleteDish_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(GetPersonalDishCountBase, nD.UserID)).WillReturnError(errors.New("database error"))

	err := repo.DeleteDish(nU.UserID, nD.PersonalDishID)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

/* I don't think we need GetUsers for anything...
func TestDb_GetUsers(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(1, "nothing@gmail.com", "Bob", "Nothing", "Bob Nothing", "2016-01-02T15:04:05",
			"ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop", false, "asdfasdfa").
		AddRow(2, "nothing2@gmail.com", "Robert", "Nothingtwo", "Robert Nothingtwo", "2016-02-02T15:04:05",
			"ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop2", false, "asdfasdfa2")

	mock.ExpectQuery(GetUsersBase).WillReturnRows(rows)

	resultingUsers, err := repo.GetUsers()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(*resultingUsers))

	resultingUser1 := (*resultingUsers)[0]
	resultingUser2 := (*resultingUsers)[1]

	assert.Equal(t, "Bob", resultingUser1.FirstName)
	assert.Equal(t, "Robert", resultingUser2.FirstName)

	assert.Equal(t, 1, resultingUser1.UserID)
	assert.Equal(t, 2, resultingUser2.UserID)
}

func TestDb_GetUsers_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"})

	mock.ExpectQuery(GetUsersBase).WillReturnRows(rows)

	resultingUsers, err := repo.GetUsers()

	assert.NotNil(t, err)
	assert.Nil(t, resultingUsers)
	//assert.Equal(t, "Database could not find any users", err.Message())
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestDb_GetUsers_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(GetUsersBase).WillReturnError(errors.New("database error"))
	resultingUsers, err := repo.GetUsers()

	assert.Nil(t, resultingUsers)
	assert.NotNil(t, err)
	//assert.Equal(t, "Error while retrieving users from the database", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetUsers_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow("SHOULDBEINT", "nothing@gmail.com", "Bob", "Nothing", "Bob Nothing", "2016-01-02T15:04:05",
			"ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop", false, "asdfasdfa")

	mock.ExpectQuery(GetUsersBase).WillReturnRows(rows)

	resultingUsers, err := repo.GetUsers()

	assert.Nil(t, resultingUsers)
	assert.NotNil(t, err)
	//assert.Equal(t, "Error while scanning the result from the database", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}
*/
func TestDb_GetUserByID(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		UserID:       2,
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-01-02T15:04:05",
		AccessToken:  "ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "105i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		Admin:        false,
		TempMatch:    "1v842d234523a",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(nU.UserID, nU.Email, nU.FirstName, nU.LastName, nU.FullName, nU.CreatedDate,
			nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.Admin, nU.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetUserByIDBase, nU.UserID)).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByID(2)

	assert.Nil(t, err)
	assert.NotNil(t, resultingUser)

	assert.Equal(t, nU.UserID, resultingUser.UserID)
	assert.Equal(t, nU.Email, resultingUser.Email)
	assert.Equal(t, nU.FirstName, resultingUser.FirstName)
	assert.Equal(t, nU.LastName, resultingUser.LastName)
	assert.Equal(t, nU.FullName, resultingUser.FullName)
	assert.Equal(t, nU.CreatedDate, resultingUser.CreatedDate)
	assert.Equal(t, nU.AccessToken, resultingUser.AccessToken)
	assert.Equal(t, nU.RefreshToken, resultingUser.RefreshToken)
	assert.Equal(t, nU.AlexaUserID, resultingUser.AlexaUserID)
	assert.Equal(t, nU.Admin, resultingUser.Admin)
	assert.Equal(t, nU.TempMatch, resultingUser.TempMatch)

}

func TestDb_GetUserByID_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(GetUserByIDBase, 1)).WillReturnError(errors.New("database error"))

	resultingUser, err := repo.GetUserByID(1)

	assert.Nil(t, resultingUser)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while retrieving user from the database", err.Message())
}

func TestDb_GetUserByID_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(GetUserByIDBase, nU.UserID)).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByID(nU.UserID)

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	assert.Equal(t, http.StatusNotFound, err.Status())
	//assert.Equal(t, "Database could not find a user with this ID", err.Message())
}

func TestDb_GetUserByID_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow("SHOULDBEINT", "nothing@gmail.com", "Bob", "Nothing", "Bob Nothing", "2016-01-02T15:04:05",
			"ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop", false, "asdfasdfa")

	mock.ExpectQuery(fmt.Sprintf(GetUserByIDBase, 1)).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByID(1)

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while scanning the result from the database", err.Message())
}

func TestDb_GetUserByID_FoundMultiple(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(1, "nothing@gmail.com", "Bob", "Nothing", "Bob Nothing", "2016-01-02T15:04:05",
			"ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop", false, "asdfasdfa").
		AddRow(2, "nothing2@gmail.com", "Robert", "Nothingtwo", "Robert Nothingtwo", "2016-02-02T15:04:05",
			"ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop2", false, "asdfasdfa2")

	mock.ExpectQuery(fmt.Sprintf(GetUserByIDBase, 1)).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByID(1)

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	//assert.Equal(t, "Database returned more than 1 row when only 1 was expected", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetUserByEmail(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		UserID:       2,
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-01-02T15:04:05",
		AccessToken:  "ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		Admin:        false,
		TempMatch:    "1v842d234523a",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(nU.UserID, nU.Email, nU.FirstName, nU.LastName, nU.FullName, nU.CreatedDate,
			nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.Admin, nU.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetUserByEmailBase, nU.Email)).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByEmail(nU.Email)

	assert.Nil(t, err)
	assert.NotNil(t, resultingUser)

	assert.Equal(t, nU.UserID, resultingUser.UserID)
	assert.Equal(t, nU.Email, resultingUser.Email)
	assert.Equal(t, nU.FirstName, resultingUser.FirstName)
	assert.Equal(t, nU.LastName, resultingUser.LastName)
	assert.Equal(t, nU.FullName, resultingUser.FullName)
	assert.Equal(t, nU.CreatedDate, resultingUser.CreatedDate)
	assert.Equal(t, nU.AccessToken, resultingUser.AccessToken)
	assert.Equal(t, nU.RefreshToken, resultingUser.RefreshToken)
	assert.Equal(t, nU.AlexaUserID, resultingUser.AlexaUserID)
	assert.Equal(t, nU.Admin, resultingUser.Admin)
	assert.Equal(t, nU.TempMatch, resultingUser.TempMatch)

}

func TestDb_GetUserByEmail_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(GetUserByEmailBase, "nothing@gmail.com")).WillReturnError(errors.New("database error"))

	resultingUser, err := repo.GetUserByEmail("nothing@gmail.com")

	assert.Nil(t, resultingUser)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while retrieving user from the database", err.Message())
}

func TestDb_GetUserByEmail_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(GetUserByEmailBase, "nothing@gmail.com")).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByEmail("nothing@gmail.com")

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	assert.Equal(t, http.StatusNotFound, err.Status())
	//assert.Equal(t, "Database could not find a user with this Email", err.Message())
}

func TestDb_GetUserByEmail_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow("SHOULDBEINT", "nothing@gmail.com", "Bob", "Nothing", "Bob Nothing", "2016-01-02T15:04:05",
			"ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop", false, "asdfasdfa")

	mock.ExpectQuery(fmt.Sprintf(GetUserByEmailBase, "nothing@gmail.com")).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByEmail("nothing@gmail.com")

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while scanning the result from the database", err.Message())
}

func TestDb_GetUserByEmail_FoundMultiple(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(1, "nothing@gmail.com", "Bob", "Nothing", "Bob Nothing", "2016-01-02T15:04:05",
			"ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop", false, "asdfasdfa").
		AddRow(2, "nothing2@gmail.com", "Robert", "Nothingtwo", "Robert Nothingtwo", "2016-02-02T15:04:05",
			"ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop2", false, "asdfasdfa2")

	mock.ExpectQuery(fmt.Sprintf(GetUserByEmailBase, "nothing@gmail.com")).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByEmail("nothing@gmail.com")

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	//assert.Equal(t, "Database returned more than 1 row when only 1 was expected", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetUserByAlexa(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		UserID:       2,
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-01-02T15:04:05",
		AccessToken:  "ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		Admin:        false,
		TempMatch:    "1v842d234523a",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(nU.UserID, nU.Email, nU.FirstName, nU.LastName, nU.FullName, nU.CreatedDate,
			nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.Admin, nU.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetUserByAlexaBase, nU.AlexaUserID)).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByAlexa(nU.AlexaUserID)

	assert.Nil(t, err)
	assert.NotNil(t, resultingUser)

	assert.Equal(t, nU.UserID, resultingUser.UserID)
	assert.Equal(t, nU.Email, resultingUser.Email)
	assert.Equal(t, nU.FirstName, resultingUser.FirstName)
	assert.Equal(t, nU.LastName, resultingUser.LastName)
	assert.Equal(t, nU.FullName, resultingUser.FullName)
	assert.Equal(t, nU.CreatedDate, resultingUser.CreatedDate)
	assert.Equal(t, nU.AccessToken, resultingUser.AccessToken)
	assert.Equal(t, nU.RefreshToken, resultingUser.RefreshToken)
	assert.Equal(t, nU.AlexaUserID, resultingUser.AlexaUserID)
	assert.Equal(t, nU.TempMatch, resultingUser.TempMatch)

}

func TestDb_GetUserByAlexa_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(GetUserByAlexaBase, "qwertyuiop")).WillReturnError(errors.New("database error"))

	resultingUser, err := repo.GetUserByAlexa("qwertyuiop")

	assert.Nil(t, resultingUser)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while retrieving user from the database", err.Message())
}

func TestDb_GetUserByAlexa_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(GetUserByAlexaBase, "qwertyuiop")).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByAlexa("qwertyuiop")

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	assert.Equal(t, http.StatusNotFound, err.Status())
	//assert.Equal(t, "Database could not find a user with this Alexa User ID", err.Message())
}

func TestDb_GetUserByAlexa_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "temp_match"}).
		AddRow("SHOULDBEINT", "nothing@gmail.com", "Bob", "Nothing", "Bob Nothing", "2016-01-02T15:04:05",
			"ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop", "asdfasdfa")

	mock.ExpectQuery(fmt.Sprintf(GetUserByAlexaBase, "qwertyuiop")).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByAlexa("qwertyuiop")

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while scanning the result from the database", err.Message())
}

func TestDb_GetUserByAlexa_FoundMultiple(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(1, "nothing@gmail.com", "Bob", "Nothing", "Bob Nothing", "2016-01-02T15:04:05",
			"ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop", false, "asdfasdfa").
		AddRow(2, "nothing2@gmail.com", "Robert", "Nothingtwo", "Robert Nothingtwo", "2016-02-02T15:04:05",
			"ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop2", false, "asdfasdfa2")

	mock.ExpectQuery(fmt.Sprintf(GetUserByAlexaBase, "qwertyuiop")).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByAlexa("qwertyuiop")

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	//assert.Equal(t, "Database returned more than 1 row when only 1 was expected", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetUserByTempMatch(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		UserID:       2,
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-01-02T15:04:05",
		AccessToken:  "ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		Admin:        false,
		TempMatch:    "1v842d234523a",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(nU.UserID, nU.Email, nU.FirstName, nU.LastName, nU.FullName, nU.CreatedDate,
			nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.Admin, nU.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetUserByTempMatchBase, nU.TempMatch)).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByTempMatch(nU.TempMatch)

	assert.Nil(t, err)
	assert.NotNil(t, resultingUser)

	assert.Equal(t, nU.UserID, resultingUser.UserID)
	assert.Equal(t, nU.Email, resultingUser.Email)
	assert.Equal(t, nU.FirstName, resultingUser.FirstName)
	assert.Equal(t, nU.LastName, resultingUser.LastName)
	assert.Equal(t, nU.FullName, resultingUser.FullName)
	assert.Equal(t, nU.CreatedDate, resultingUser.CreatedDate)
	assert.Equal(t, nU.AccessToken, resultingUser.AccessToken)
	assert.Equal(t, nU.RefreshToken, resultingUser.RefreshToken)
	assert.Equal(t, nU.AlexaUserID, resultingUser.AlexaUserID)
	assert.Equal(t, nU.Admin, resultingUser.Admin)
	assert.Equal(t, nU.TempMatch, resultingUser.TempMatch)

}

func TestDb_GetUserByTempMatch_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(GetUserByTempMatchBase, "qwertyuiop")).WillReturnError(errors.New("database error"))

	resultingUser, err := repo.GetUserByTempMatch("qwertyuiop")

	assert.Nil(t, resultingUser)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while retrieving user from the database", err.Message())
}

func TestDb_GetUserByTempMatch_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(GetUserByTempMatchBase, "qwertyuiop")).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByTempMatch("qwertyuiop")

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	assert.Equal(t, http.StatusNotFound, err.Status())
	//assert.Equal(t, "Database could not find a user with this Temp Match", err.Message())
}

func TestDb_GetUserByTempMatch_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow("SHOULDBEINT", "nothing@gmail.com", "Bob", "Nothing", "Bob Nothing", "2016-01-02T15:04:05",
			"ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop", false, "asdfasdfa")

	mock.ExpectQuery(fmt.Sprintf(GetUserByTempMatchBase, "qwertyuiop")).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByTempMatch("qwertyuiop")

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while scanning the result from the database", err.Message())
}

func TestDb_GetUserByTempMatch_FoundMultiple(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(1, "nothing@gmail.com", "Bob", "Nothing", "Bob Nothing", "2016-01-02T15:04:05",
			"ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "1//05i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop", false, "asdfasdfa").
		AddRow(2, "nothing2@gmail.com", "Robert", "Nothingtwo", "Robert Nothingtwo", "2016-02-02T15:04:05",
			"ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop2", false, "asdfasdfa2")

	mock.ExpectQuery(fmt.Sprintf(GetUserByTempMatchBase, "qwertyuiop")).WillReturnRows(rows)

	resultingUser, err := repo.GetUserByTempMatch("qwertyuiop")

	assert.NotNil(t, err)
	assert.Nil(t, resultingUser)

	//assert.Equal(t, "Database returned more than 1 row when only 1 was expected", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_CreateUser(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-02-02T15:04:05",
		AccessToken:  "ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		Admin:        false,
	}

	createRows := sqlmock.NewRows([]string{""})

	getRows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(1, "nothing@gmail.com", "Bob", "Nothing", "Bob Nothing", "2016-01-02T15:04:05",
			"ya33.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k", "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM", "qwertyuiop", false, "adfasfsgas654g")

	mock.ExpectQuery(`INSERT INTO user \(.+\) VALUES\(".+", "Bob", "Nothing", ".+", ".+", ".+", ".+", "qwertyuiop", false, ".+"\)`).
		WillReturnRows(createRows)

	mock.ExpectQuery(`SELECT \* FROM user WHERE temp_match = ".+"`).WillReturnRows(getRows)

	returnedUser, err := repo.CreateUser(*nU)

	assert.Nil(t, err)

	assert.NotNil(t, returnedUser)
	assert.Equal(t, nU.FirstName, returnedUser.FirstName)
}

func TestDb_CreateUser_InsertError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-02-02T15:04:05",
		AccessToken:  "ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		Admin:        false,
		TempMatch:    "a4s65df6adhy4s5gjet",
	}

	mock.ExpectQuery(`INSERT INTO user \(.+\) VALUES\(".+", "Bob", "Nothing", ".+", ".+", ".+", ".+", "qwertyuiop", false, ".+"\)`).
		WillReturnError(errors.New("not possible"))

	returnedUser, err := repo.CreateUser(*nU)

	assert.NotNil(t, err)
	assert.Nil(t, returnedUser)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while inserting the user into the database", err.Message())
}

func TestDb_CreateUser_CheckError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-02-02T15:04:05",
		AccessToken:  "ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		Admin:        false,
		TempMatch:    "a4s65df6adhy4s5gjet",
	}

	createRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(`INSERT INTO user \(.+\) VALUES\(".+", "Bob", "Nothing", ".+", ".+", ".+", ".+", "qwertyuiop", false, ".+"\)`).
		WillReturnRows(createRows)

	mock.ExpectQuery(`SELECT \* FROM user WHERE temp_match = ".+"`).
		WillReturnError(errors.New("not possible"))

	returnedUser, err := repo.CreateUser(*nU)

	assert.NotNil(t, err)
	assert.Nil(t, returnedUser)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while checking the user that was created."+
	//	" Cannot verify if anything was entered to the Database", err.Message())
}

func TestDb_UpdateUser(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		UserID:       2,
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-02-02T15:04:05",
		AccessToken:  "ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		Admin:        false,
		TempMatch:    "a4s65df6adhy4s5gjet",
	}

	updateRows := sqlmock.NewRows([]string{""})

	getRows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(nU.UserID, nU.Email, nU.FirstName, nU.LastName, nU.FullName, nU.CreatedDate,
			nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.Admin, nU.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(UpdateUserBase, nU.Email, nU.FirstName, nU.LastName, nU.FullName,
		nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.TempMatch, nU.UserID)).WillReturnRows(updateRows)

	mock.ExpectQuery(fmt.Sprintf(GetUserByIDBase, nU.UserID)).WillReturnRows(getRows)

	returnedUser, err := repo.UpdateUser(*nU)

	assert.Nil(t, err)
	assert.NotNil(t, returnedUser)

	assert.Equal(t, nU.UserID, returnedUser.UserID)
	assert.Equal(t, nU.Email, returnedUser.Email)
	assert.Equal(t, nU.FirstName, returnedUser.FirstName)
	assert.Equal(t, nU.LastName, returnedUser.LastName)
	assert.Equal(t, nU.FullName, returnedUser.FullName)
	assert.Equal(t, nU.CreatedDate, returnedUser.CreatedDate)
	assert.Equal(t, nU.AccessToken, returnedUser.AccessToken)
	assert.Equal(t, nU.RefreshToken, returnedUser.RefreshToken)
	assert.Equal(t, nU.AlexaUserID, returnedUser.AlexaUserID)
	assert.Equal(t, nU.Admin, returnedUser.Admin)
	assert.Equal(t, nU.TempMatch, returnedUser.TempMatch)

}

func TestDb_UpdateUser_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		UserID:       2,
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-02-02T15:04:05",
		AccessToken:  "ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		TempMatch:    "a4s65df6adhy4s5gjet",
	}

	mock.ExpectQuery(fmt.Sprintf(UpdateUserBase, nU.Email, nU.FirstName, nU.LastName, nU.FullName,
		nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.TempMatch, nU.UserID)).
		WillReturnError(errors.New("database error"))

	returnedUser, err := repo.UpdateUser(*nU)

	assert.Nil(t, returnedUser)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while updating the user in the database", err.Message())
}

func TestDb_UpdateUser_CheckError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		UserID:       2,
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-02-02T15:04:05",
		AccessToken:  "ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		TempMatch:    "a4s65df6adhy4s5gjet",
	}

	createRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(fmt.Sprintf(UpdateUserBase, nU.Email, nU.FirstName, nU.LastName, nU.FullName,
		nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.TempMatch, nU.UserID)).
		WillReturnRows(createRows)

	mock.ExpectQuery(fmt.Sprintf(GetUserByIDBase, nU.UserID)).WillReturnError(errors.New("database error"))

	returnedUser, err := repo.UpdateUser(*nU)

	assert.Nil(t, returnedUser)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while checking the user that was created."+
	//	" Cannot verify if anything was updated in the Database", err.Message())

}

func TestDb_DeleteUser(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		UserID:       2,
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-02-02T15:04:05",
		AccessToken:  "ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		TempMatch:    "a4s65df6adhy4s5gjet",
	}

	deleteRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(fmt.Sprintf(DeleteUserBase, nU.UserID)).WillReturnRows(deleteRows)

	err := repo.DeleteUser(*nU)

	assert.Nil(t, err)
}

func TestDb_DeleteUser_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		UserID:       2,
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-02-02T15:04:05",
		AccessToken:  "ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		TempMatch:    "a4s65df6adhy4s5gjet",
	}

	mock.ExpectQuery(fmt.Sprintf(DeleteUserBase, nU.UserID)).WillReturnError(errors.New("database error"))

	err := repo.DeleteUser(*nU)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while deleting the user from the database", err.Message())
}

func TestDb_DeleteUser_CheckError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	nU := &user.User{
		UserID:       2,
		Email:        "nothing@gmail.com",
		FirstName:    "Bob",
		LastName:     "Nothing",
		FullName:     "Bob Nothing",
		CreatedDate:  "2016-02-02T15:04:05",
		AccessToken:  "ya44.a0Ae4lvC1iHeKSDRdQ542I-lEy8LHUU7-9r-k",
		RefreshToken: "205i7nDY0JDTJmCgYIAQDKJSNwF-L9IrRgJ4-fM",
		AlexaUserID:  "qwertyuiop",
		Admin:        false,
		TempMatch:    "a4s65df6adhy4s5gjet",
	}

	getRows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(nU.UserID, nU.Email, nU.FirstName, nU.LastName, nU.FullName, nU.CreatedDate,
			nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.Admin, nU.TempMatch)

	deleteRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(fmt.Sprintf(DeleteUserBase, nU.UserID)).WillReturnRows(deleteRows)

	mock.ExpectQuery(fmt.Sprintf(GetUserByIDBase, nU.UserID)).WillReturnRows(getRows)

	err := repo.DeleteUser(*nU)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while deleting the user from the database, could not verify it was deleted.", err.Message())
}

func TestDb_GetStoragesByUser(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "title", "description", "temp_match"}).
		AddRow(nS.StorageID, nS.PersonalID, nS.UserID, nS.Title, nS.Description, nS.TempMatch).
		AddRow(nS.StorageID+1, nS.PersonalID+1, nS.UserID, nS.Title+"2", nS.Description+"2", nS.TempMatch+"2")

	mock.ExpectQuery(fmt.Sprintf(GetStoragesBase, nS.UserID)).WillReturnRows(rows)

	resultingStorages, err := repo.GetStorages(nS.UserID)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(*resultingStorages))

	resultingStorage1 := (*resultingStorages)[0]
	resultingStorage2 := (*resultingStorages)[1]

	assert.Equal(t, *nS, resultingStorage1)
	assert.NotEqual(t, resultingStorage1, resultingStorage2)
}

func TestDb_GetStoragesByUser_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "title", "description", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(GetStoragesBase, nS.UserID)).WillReturnRows(rows)

	resultingStorages, err := repo.GetStorages(nS.UserID)

	assert.NotNil(t, err)
	assert.Nil(t, resultingStorages)
	//assert.Equal(t, "Database could not find any storage units for this user", err.Message())
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestDb_GetStoragesByUser_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(GetStoragesBase, 1)).WillReturnError(errors.New("database error"))

	resultingStorages, err := repo.GetStorages(nS.UserID)

	assert.Nil(t, resultingStorages)
	assert.NotNil(t, err)
	//assert.Equal(t, "Error while retrieving storage units from the database", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetStoragesByUser_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "title", "description", "temp_match"}).
		AddRow("SHOULD BE INT", nS.PersonalID, nS.UserID, nS.Title, nS.Description, nS.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetStoragesBase, nS.UserID)).WillReturnRows(rows)

	resultingStorages, err := repo.GetStorages(nS.UserID)

	assert.Nil(t, resultingStorages)
	assert.NotNil(t, err)
	//assert.Equal(t, "Error while scanning the result from the database", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetStorageByID(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "title", "description", "temp_match"}).
		AddRow(nS.StorageID, nS.PersonalID, nS.UserID, nS.Title, nS.Description, nS.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetStorageByIDBase, nS.UserID, nS.PersonalID)).WillReturnRows(rows)

	resultingStorage, err := repo.GetStorageByID(nS.UserID, nS.PersonalID)

	assert.Nil(t, err)
	assert.NotNil(t, resultingStorage)

	assert.Equal(t, nS, resultingStorage)

}

func TestDb_GetStorageByID_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(GetStorageByIDBase, nS.UserID, nS.PersonalID)).WillReturnError(errors.New("database error"))

	resultingStorage, err := repo.GetStorageByID(nS.UserID, nS.PersonalID)

	assert.Nil(t, resultingStorage)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while retrieving storage unit from the database", err.Message())
}

func TestDb_GetStorageByID_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "title", "description", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(GetStorageByIDBase, nS.UserID, nS.PersonalID)).WillReturnRows(rows)

	resultingStorage, err := repo.GetStorageByID(nS.UserID, nS.PersonalID)

	assert.NotNil(t, err)
	assert.Nil(t, resultingStorage)

	assert.Equal(t, http.StatusNotFound, err.Status())
	//assert.Equal(t, "Database could not find a storage unit with this ID", err.Message())
}

func TestDb_GetStorageByID_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "title", "description", "temp_match"}).
		AddRow(nS.StorageID, "SHOULD BE INT", nS.UserID, nS.Title, nS.Description, nS.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetStorageByIDBase, nS.UserID, nS.PersonalID)).WillReturnRows(rows)

	resultingStorage, err := repo.GetStorageByID(nS.UserID, nS.PersonalID)

	assert.NotNil(t, err)
	assert.Nil(t, resultingStorage)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while scanning the result from the database", err.Message())
}

func TestDb_GetStorageByID_FoundMultiple(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "title", "description", "temp_match"}).
		AddRow(nS.StorageID, nS.PersonalID, nS.UserID, nS.Title, nS.Description, nS.TempMatch).
		AddRow(nS.StorageID+1, nS.PersonalID+1, nS.UserID, nS.Title+"2", nS.Description+"2", nS.TempMatch+"2")

	mock.ExpectQuery(fmt.Sprintf(GetStorageByIDBase, nS.UserID, nS.PersonalID)).WillReturnRows(rows)

	resultingStorage, err := repo.GetStorageByID(nS.UserID, nS.PersonalID)

	assert.NotNil(t, err)
	assert.Nil(t, resultingStorage)

	//assert.Equal(t, "Database returned more than 1 row when only 1 was expected", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_CreateStorage(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	createRows := sqlmock.NewRows([]string{""})

	getRows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "title", "description", "temp_match"}).
		AddRow(nS.StorageID, nS.PersonalID, nS.UserID, nS.Title, nS.Description, nS.TempMatch)

	mock.ExpectQuery(`INSERT INTO storage \(personal_id, user_id, title, description, temp_match\) VALUES\(.+\)`).
		WillReturnRows(createRows)

	mock.ExpectQuery(`SELECT \* FROM storage WHERE temp_match=".+"`).WillReturnRows(getRows)

	returnedStorage, err := repo.CreateStorage(*nS)

	assert.Nil(t, err)

	assert.NotNil(t, returnedStorage)
	assert.Equal(t, nS.Title, returnedStorage.Title)
}

func TestDb_CreateStorage_InsertError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(`INSERT INTO storage \(user_id, title, description, temp_match\) VALUES\(.+\)`).
		WillReturnError(errors.New("not possible"))

	returnedStorage, err := repo.CreateStorage(*nS)

	assert.NotNil(t, err)
	assert.Nil(t, returnedStorage)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while inserting the storage unit into the database", err.Message())
}

func TestDb_CreateStorage_CheckError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	createRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(`INSERT INTO storage \(user_id, title, description, temp_match\) VALUES\(.+\)`).
		WillReturnRows(createRows)

	mock.ExpectQuery(`SELECT \* FROM storage WHERE temp_match=".+"`).WillReturnError(errors.New("database error"))

	returnedStorage, err := repo.CreateStorage(*nS)

	assert.NotNil(t, err)
	assert.Nil(t, returnedStorage)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while checking the storage unit that was created."+
	//" Cannot verify if anything was entered to the Database", err.Message())
}

func TestDb_UpdateStorage(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	updateRows := sqlmock.NewRows([]string{""})

	getRows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "title", "description", "temp_match"}).
		AddRow(nS.StorageID, nS.PersonalID, nS.UserID, nS.Title, nS.Description, nS.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(UpdateStorageBase, nS.PersonalID, nS.Title, nS.Description, nS.TempMatch, nS.StorageID)).WillReturnRows(updateRows)

	mock.ExpectQuery(fmt.Sprintf(GetStorageByIDBase, nS.UserID, nS.PersonalID)).WillReturnRows(getRows)

	err := repo.UpdateStorage(*nS)

	assert.Nil(t, err)

}

func TestDb_UpdateStorage_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(UpdateStorageBase, nS.PersonalID, nS.Title, nS.Description, nS.TempMatch, nS.StorageID)).
		WillReturnError(errors.New("database error"))

	err := repo.UpdateStorage(*nS)

	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while updating the storage unit in the database", err.Message())
}

func TestDb_UpdateStorage_CheckError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	updateRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(fmt.Sprintf(UpdateStorageBase, nS.PersonalID, nS.Title, nS.Description, nS.TempMatch, nS.StorageID)).WillReturnRows(updateRows)

	mock.ExpectQuery(fmt.Sprintf(GetStorageByIDBase, nS.UserID, nS.PersonalID)).WillReturnError(errors.New("database error"))

	err := repo.UpdateStorage(*nS)

	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while checking the storage unit that was created."+
	//	" Cannot verify if anything was updated in the Database", err.Message())

}

func TestDb_DeleteStorage(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	deleteRows := sqlmock.NewRows([]string{""})

	mock.ExpectQuery(fmt.Sprintf(DeleteStorageBase, nS.UserID, nS.PersonalID)).WillReturnRows(deleteRows)

	err := repo.DeleteStorage(nS.UserID, nS.PersonalID)

	assert.Nil(t, err)
}

func TestDb_DeleteStorage_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(DeleteStorageBase, nS.UserID, nS.PersonalID)).WillReturnError(errors.New("database error"))

	err := repo.DeleteStorage(nS.UserID, nS.PersonalID)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while deleting the storage unit from the database", err.Message())
}

func TestDb_DeleteStorage_CheckError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	deleteRows := sqlmock.NewRows([]string{""})

	getRows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "title", "description", "temp_match"}).
		AddRow(nS.StorageID, nS.PersonalID, nS.UserID, nS.Title, nS.Description, nS.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(DeleteStorageBase, nS.UserID, nS.PersonalID)).WillReturnRows(deleteRows)

	mock.ExpectQuery(fmt.Sprintf(GetStorageByIDBase, nS.UserID, nS.PersonalID)).WillReturnRows(getRows)

	err := repo.DeleteStorage(nS.UserID, nS.PersonalID)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
	//assert.Equal(t, "Error while deleting the storage unit from the database, could not verify it was deleted.", err.Message())
}

func TestDb_GetStorageDishes(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description,
			nD.CreatedDate, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch).
		AddRow(nD.DishID+200, nD.PersonalDishID+1, nD.UserID, nD.StorageID, nD.Title+"2", nD.Description+"2",
			nD.CreatedDate, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch+"2")

	mock.ExpectQuery(fmt.Sprintf(GetStorageDishesBase, nS.UserID, nS.PersonalID)).WillReturnRows(rows)

	resultingDishes, err := repo.GetStorageDishes(nS.UserID, nS.PersonalID)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(*resultingDishes))

	resultingDish1 := (*resultingDishes)[0]
	resultingDish2 := (*resultingDishes)[1]

	assert.Equal(t, *nD, resultingDish1)
	assert.NotEqual(t, resultingDish1, resultingDish2)
}

func TestDb_GetStorageDishes_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(GetStorageDishesBase, nS.UserID, nS.PersonalID)).WillReturnRows(rows)

	resultingDishes, err := repo.GetStorageDishes(nS.UserID, nS.PersonalID)

	assert.NotNil(t, err)
	assert.Nil(t, resultingDishes)
	//assert.Equal(t, "Database could not find any dishes that belong to this storage unit", err.Message())
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestDb_GetStorageDishes_QueryError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(fmt.Sprintf(GetStorageDishesBase, nS.UserID, nS.PersonalID)).WillReturnError(errors.New("database error"))
	resultingDishes, err := repo.GetStorageDishes(nS.UserID, nS.PersonalID)

	assert.Nil(t, resultingDishes)
	assert.NotNil(t, err)
	//assert.Equal(t, "Error while retrieving dishes from the database", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDb_GetStorageDishes_RowScanError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow("SHOULD BE INT", nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description,
			nD.CreatedDate, nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(GetStorageDishesBase, nS.UserID, nS.PersonalID)).WillReturnRows(rows)

	resultingDishes, err := repo.GetStorageDishes(nS.UserID, nS.PersonalID)

	assert.Nil(t, resultingDishes)
	assert.NotNil(t, err)
	//assert.Equal(t, "Error while scanning the result from the database", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}
