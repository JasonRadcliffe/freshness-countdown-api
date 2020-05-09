package dish

import (
	"fmt"
	"testing"

	dishDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	userDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/user"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	dbrepo "github.com/jasonradcliffe/freshness-countdown-api/repository/db"
	"github.com/stretchr/testify/assert"
)

var nD = &dishDomain.Dish{
	DishID:         1,
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

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishesBase)).WillReturnRows(rows)

	resultingDishes, err := dS.GetAll(nU)
	dish := (*resultingDishes)[0]

	assert.Nil(t, err)
	assert.NotNil(t, resultingDishes)
	assert.Equal(t, nD, &dish)
	assert.Equal(t, 2, len(*resultingDishes))

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

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishesBase)).WillReturnRows(rows)

	resultingDishes, err := dS.GetExpired(nU)
	//dish := (*resultingDishes)[0]

	assert.Nil(t, err)
	assert.NotNil(t, resultingDishes)
	assert.Equal(t, 1, len(*resultingDishes))

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
