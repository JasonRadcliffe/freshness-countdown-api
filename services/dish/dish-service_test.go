package dish

import (
	"fmt"
	"testing"

	dishDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/dish"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	dbrepo "github.com/jasonradcliffe/freshness-countdown-api/repository/db"
	"github.com/stretchr/testify/assert"
)

var nD = &dishDomain.Dish{
	DishID:      1,
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

	rows := sqlmock.NewRows([]string{"id", "user_id", "storage_id", "title", "description", "created_date",
		"expire_date", "priority", "dish_type", "portions", "temp_match"}).
		AddRow(nD.DishID, nD.UserID, nD.StorageID, nD.Title, nD.Description, nD.CreatedDate,
			nD.ExpireDate, nD.Priority, nD.DishType, nD.Portions, nD.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetDishByIDBase, 1)).WillReturnRows(rows)

	resultingDish, err := dS.GetByID(1)
	fmt.Println("got this dish from the test:", resultingDish)

	assert.Equal(t, nD.Title, resultingDish.Title)
}
