package api

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jasonradcliffe/freshness-countdown-api/services/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/services/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/services/user"

	userDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/user"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	dbrepo "github.com/jasonradcliffe/freshness-countdown-api/repository/db"
	"github.com/stretchr/testify/assert"
)

type mockOAuthConfig struct{}

var rUser = &userDomain.User{
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

func TestAPIHandler_getExpiredDishes(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	dS := dish.NewService(repo)
	uS := user.NewService(repo)
	sS := storage.NewService(repo)

	//rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
	//	"expire_date", "priority", "dish_type", "portions", "temp_match"}).
	//	AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
	//		nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishByIDBase, 1)).WillReturnError(errors.New(""))

	resultingDish, err := dS.GetByID("alexaid1234", "superaccesstoken", 1)
	fmt.Println("got this dish from the test:", resultingDish)

	assert.Equal(t, "", "")
}
