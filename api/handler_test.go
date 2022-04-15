package api

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"golang.org/x/oauth2"

	"github.com/jasonradcliffe/freshness-countdown-api/services/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/services/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/services/user"

	dishDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
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

var nDex = &dishDomain.Dish{
	DishID:         1,
	PersonalDishID: 1,
	UserID:         2,
	StorageID:      3,
	Title:          "Carrots",
	Description:    "Some carrots we got at the store",
	CreatedDate:    "2006-01-02T15:04:05",
	ExpireDate:     "2019-10-13T08:00",
	Priority:       "",
	DishType:       "",
	Portions:       -1,
	TempMatch:      "9r842d3a351",
}

//Exchange is the mock method to get the token for oauth
func (m *mockOAuthConfig) AuthCodeURL(state string, options ...oauth2.AuthCodeOption) string {
	return ""
}

//Exchange is the mock method to get the token for oauth
func (m *mockOAuthConfig) Exchange(c context.Context, code string, options ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return nil, nil
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

	oC := &mockOAuthConfig{}

	dS := dish.NewService(repo)
	uS := user.NewService(repo)
	sS := storage.NewService(repo)

	mHandler := NewHandler(dS, sS, uS, oC)
	fmt.Println("testing:", mHandler)

	rows := sqlmock.NewRows([]string{"id", "personal_id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nDex.DishID, nDex.PersonalDishID, nDex.UserID, nDex.StorageID, nDex.Title, nDex.Description, nDex.CreatedDate,
			nDex.ExpireDate, nDex.Priority, nDex.DishType, nDex.Portions, nDex.TempMatch).
		AddRow(nD.DishID, nD.PersonalDishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishesBase, rUser.UserID)).WillReturnRows(rows)

	resultingDishesMarshaled, err := getDishesExpired(rUser, dS)
	var resultingDishes dishDomain.Dishes
	jsonErr := json.Unmarshal(resultingDishesMarshaled, &resultingDishes)

	assert.Nil(t, jsonErr)

	assert.Equal(t, 2, len(resultingDishes))
}
