package dish

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	dishDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	userDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/user"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	dbrepo "github.com/jasonradcliffe/freshness-countdown-api/repository/db"
	"github.com/stretchr/testify/assert"
)

var nD = &dishDomain.Dish{
	DishID:         200,
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

var nU = &userDomain.User{
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

func TestDishService_GetByID(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishByIDBase, nD.UserID, nD.PersonalDishID)).WillReturnRows(rows)

	resultingDish, err := dS.GetByID(nU, nD.PersonalDishID)
	fmt.Println("got this dish from the test:", resultingDish)

	assert.Equal(t, nD.Title, resultingDish.Title)
}

func TestDishService_GetByID_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishByIDBase, nD.UserID, nD.PersonalDishID)).WillReturnRows(rows)

	resultingDish, err := dS.GetByID(nU, nD.PersonalDishID)
	fmt.Println("got this dish from the test:", resultingDish)

	assert.Nil(t, resultingDish)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDishService_GetAll(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch).
		AddRow(nD.DishID+1, nD.PersonalDishID+1, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishesBase, nU.UserID)).WillReturnRows(rows)

	resultingDishes, err := dS.GetAll(nU)
	dish := (*resultingDishes)[0]

	assert.Nil(t, err)
	assert.NotNil(t, resultingDishes)
	assert.Equal(t, nD, &dish)
	assert.Equal(t, 2, len(*resultingDishes))

}

func TestDishService_GetAll_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishesBase, nU.UserID)).WillReturnRows(rows)

	resultingDishes, err := dS.GetAll(nU)

	assert.Nil(t, resultingDishes)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())

}

func TestDishService_GetExpired(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch).
		AddRow(nD.DishID+1, nD.PersonalDishID+1, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			"2019-10-13T08:00", nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishesBase, nU.UserID)).WillReturnRows(rows)

	resultingDishes, err := dS.GetExpired(nU)
	//dish := (*resultingDishes)[0]

	assert.Nil(t, err)
	assert.NotNil(t, resultingDishes)
	assert.Equal(t, 2, len(*resultingDishes))

}

func TestDishService_GetExpired_InvalidDishExpDateInList(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch).
		AddRow(nD.DishID+1, nD.PersonalDishID+1, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			"201910INVALIDDATE13T08:00", nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishesBase, nU.UserID)).WillReturnRows(rows)

	resultingDishes, err := dS.GetExpired(nU)
	//dish := (*resultingDishes)[0]

	assert.Nil(t, err)
	assert.NotNil(t, resultingDishes)
	assert.Equal(t, 1, len(*resultingDishes))

}

func TestDishService_GetExpired_NoDishesFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishesBase, nU.UserID)).WillReturnRows(rows)

	resultingDishes, err := dS.GetExpired(nU)
	//dish := (*resultingDishes)[0]

	assert.Nil(t, resultingDishes)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())

}

func TestDishService_GetExpiredByDate(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch).
		AddRow(nD.DishID+1, nD.PersonalDishID+1, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			"2024INVALID10-13T08:00", nD.Priority, nD.DishType, nD.Portions, nD.TempMatch).
		AddRow(nD.DishID+1, nD.PersonalDishID+1, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			"2024-10-13T08:00", nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishesBase, nU.UserID)).WillReturnRows(rows)

	resultingDishes, err := dS.GetExpiredByDate(nU, "2023-10-13T08:00")
	dish := (*resultingDishes)[0]

	assert.Nil(t, err)
	assert.NotNil(t, resultingDishes)
	assert.Equal(t, nD, &dish)
	assert.Equal(t, 1, len(*resultingDishes))

}

func TestDishService_GetExpiredByDate_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishesBase, nU.UserID)).WillReturnRows(rows)

	resultingDishes, err := dS.GetExpiredByDate(nU, nD.ExpireDate)

	assert.Nil(t, resultingDishes)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())

}

func TestDishService_Create(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)

	mock.ExpectQuery(`I.*`).WillReturnRows(rows)

	mock.ExpectQuery(`SELECT \* FROM dish WHERE temp_match = ".+"`).WillReturnRows(rows)

	resultingDish, err := dS.Create(nU, nD, "P1Y3DT2M")

	assert.Nil(t, err)
	assert.NotNil(t, resultingDish)

}

func TestDishService_Create2MoreTimeCombinations(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)

	mock.ExpectQuery(`I.*`).WillReturnRows(rows)

	mock.ExpectQuery(`SELECT \* FROM dish WHERE temp_match = ".+"`).WillReturnRows(rows)

	resultingDish, err := dS.Create(nU, nD, "P1MT2H30S")

	assert.Nil(t, err)
	assert.NotNil(t, resultingDish)

}

func TestDishService_Create2MoreTimeCombinations_ParseErrors1(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)
	mock.ExpectQuery(`I.*`).WillReturnRows(rows)
	mock.ExpectQuery(`SELECT \* FROM dish WHERE temp_match = ".+"`).WillReturnRows(rows)

	resultingDish, err := dS.Create(nU, nD, "P1aY1M1DT2H2M30S")
	assert.Nil(t, err)
	assert.NotNil(t, resultingDish)

}

func TestDishService_Create2MoreTimeCombinations_ParseErrors2(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)
	mock.ExpectQuery(`I.*`).WillReturnRows(rows)
	mock.ExpectQuery(`SELECT \* FROM dish WHERE temp_match = ".+"`).WillReturnRows(rows)

	resultingDish, err := dS.Create(nU, nD, "P1Ya1M1DT2H2M30S")
	assert.Nil(t, err)
	assert.NotNil(t, resultingDish)

}

func TestDishService_Create2MoreTimeCombinations_ParseErrors3(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)
	mock.ExpectQuery(`I.*`).WillReturnRows(rows)
	mock.ExpectQuery(`SELECT \* FROM dish WHERE temp_match = ".+"`).WillReturnRows(rows)

	resultingDish, err := dS.Create(nU, nD, "P1Y1M1bDT2H2M30S")
	assert.Nil(t, err)
	assert.NotNil(t, resultingDish)

}

func TestDishService_Create2MoreTimeCombinations_ParseErrors4(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)
	mock.ExpectQuery(`I.*`).WillReturnRows(rows)
	mock.ExpectQuery(`SELECT \* FROM dish WHERE temp_match = ".+"`).WillReturnRows(rows)

	resultingDish, err := dS.Create(nU, nD, "P1Y1M1DTf2H2M30S")
	assert.Nil(t, err)
	assert.NotNil(t, resultingDish)

}

func TestDishService_Create2MoreTimeCombinations_ParseErrors5(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)
	mock.ExpectQuery(`I.*`).WillReturnRows(rows)
	mock.ExpectQuery(`SELECT \* FROM dish WHERE temp_match = ".+"`).WillReturnRows(rows)

	resultingDish, err := dS.Create(nU, nD, "P1Y1M1DT2Hn2M30S")
	assert.Nil(t, err)
	assert.NotNil(t, resultingDish)

}

func TestDishService_Create2MoreTimeCombinations_ParseErrors6(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)
	mock.ExpectQuery(`I.*`).WillReturnRows(rows)
	mock.ExpectQuery(`SELECT \* FROM dish WHERE temp_match = ".+"`).WillReturnRows(rows)

	resultingDish, err := dS.Create(nU, nD, "P1Y1M1DT2H2M3d0S")
	assert.Nil(t, err)
	assert.NotNil(t, resultingDish)

}

func TestDishService_Create_ErrorOnDishCountLookup(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	mock.ExpectQuery(`SELECT C.*`).WillReturnError(errors.New("database could not perform this action or returned some error."))

	resultingDish, err := dS.Create(nU, nD, "P2DT2H")

	assert.Nil(t, resultingDish)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())

}

func TestDishService_Create_CouldNotCreate(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)

	mock.ExpectQuery(`I.*`).WillReturnRows(rows)

	mock.ExpectQuery(`SELECT \* FROM dish WHERE temp_match = ".+"`).WillReturnError(errors.New("Could Not Retrieve a dish with the same temp match as we though we just added"))

	resultingDish, err := dS.Create(nU, nD, "P2DT2H")

	assert.Nil(t, resultingDish)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())

}

func TestDishService_Update(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	emptyRows := sqlmock.NewRows([]string{})

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`UPDATE.*`).WillReturnRows(emptyRows)

	mock.ExpectQuery(`SELECT \* FROM dish WHERE.*`).WillReturnRows(rows)

	err = dS.Update(nU, nD, "P1Y3DT2M")

	assert.Nil(t, err)
}

func TestDishService_Update_CouldNotUpdate(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	mock.ExpectQuery(`UPDATE.*`).WillReturnError(errors.New("Database error, could not update"))

	err = dS.Update(nU, nD, "P1Y3DT2M")

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())

}

func TestDishService_Update_CouldNotVerifyUpdate(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	emptyRows := sqlmock.NewRows([]string{})

	mock.ExpectQuery(`UPDATE.*`).WillReturnRows(emptyRows)

	mock.ExpectQuery(`SELECT \* FROM dish WHERE.*`).WillReturnError(errors.New("Database error - could not verify update"))

	err = dS.Update(nU, nD, "P1Y3DT2M")

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDishService_Delete(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)

	emptyRows := sqlmock.NewRows([]string{})

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)

	mock.ExpectQuery(`DELETE FROM dish WHERE.*`).WillReturnRows(emptyRows)

	mock.ExpectQuery(`SELECT \* FROM dish WHERE.*`).WillReturnError(errors.New("Database error - dish not found"))

	err = dS.Delete(nU, nD.PersonalDishID)

	assert.Nil(t, err)
}

func TestDishService_Delete_DishIDTooHigh(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)

	emptyRows := sqlmock.NewRows([]string{})

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)

	mock.ExpectQuery(`DELETE FROM dish WHERE.*`).WillReturnRows(emptyRows)

	mock.ExpectQuery(`SELECT \* FROM dish WHERE.*`).WillReturnError(errors.New("Database error - dish not found"))

	err = dS.Delete(nU, nD.PersonalDishID+2)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status())
}

func TestDishService_Delete_ErrorGettingDishCount(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	mock.ExpectQuery(`SELECT C.*`).WillReturnError(errors.New("Database error - could not get dish count"))

	err = dS.Delete(nU, nD.PersonalDishID)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDishService_Delete_ErrorOnDelete(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)

	mock.ExpectQuery(`DELETE FROM dish WHERE.*`).WillReturnError(errors.New("Could not do the delete query"))

	err = dS.Delete(nU, nD.PersonalDishID)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDishService_Delete_FindsDeletedDishOnDoubleCheck(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)

	emptyRows := sqlmock.NewRows([]string{})

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)

	mock.ExpectQuery(`DELETE FROM dish WHERE.*`).WillReturnRows(emptyRows)

	mock.ExpectQuery(`SELECT \* FROM dish WHERE.*`).WillReturnRows(rows)

	err = dS.Delete(nU, nD.PersonalDishID)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestDishService_Delete_DecrementSomeDishes(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(4)

	emptyRows := sqlmock.NewRows([]string{})

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)

	mock.ExpectQuery(`DELETE FROM dish WHERE.*`).WillReturnRows(emptyRows)

	mock.ExpectQuery(`SELECT \* FROM dish WHERE.*`).WillReturnError(errors.New("Database error - dish not found"))

	mock.ExpectQuery(`UPDATE.*`).WillReturnRows(emptyRows)

	err = dS.Delete(nU, nD.PersonalDishID)

	assert.Nil(t, err)
}

func TestDishService_Delete_DecrementSomeDishesError(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := NewService(repo)

	dishCount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(4)

	emptyRows := sqlmock.NewRows([]string{})

	mock.ExpectQuery(`SELECT C.*`).WillReturnRows(dishCount)

	mock.ExpectQuery(`DELETE FROM dish WHERE.*`).WillReturnRows(emptyRows)

	mock.ExpectQuery(`SELECT \* FROM dish WHERE.*`).WillReturnError(errors.New("Database error - dish not found"))

	mock.ExpectQuery(`UPDATE.*`).WillReturnError(errors.New("Could not decrement those dishes"))

	err = dS.Delete(nU, nD.PersonalDishID)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())

}
