package user

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	userDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/user"
	dbrepo "github.com/jasonradcliffe/freshness-countdown-api/repository/db"

	"github.com/stretchr/testify/assert"
)

const (
	googleAPIOKResponse = `{
		"sub": "114668774842776472919",
		"name": "Bob Nothing",
		"given_name": "Bob",
		"family_name": "Nothing",
		"picture": "https://lh3.googleusercontent.com/a-/AOh14GjaNZnU1_PuxYIb9tO_3uVMV3e",
		"email": "nothing@gmail.com",
		"email_verified": true,
		"locale": "en"
	  }`
)

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

func testHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewTLSServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return cli, s.Close
}

func TestUser_GetByID(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	userService := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(nU.UserID, nU.Email, nU.FirstName, nU.LastName, nU.FullName, nU.CreatedDate,
			nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.Admin, nU.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetUserByIDBase, nU.UserID)).WillReturnRows(rows)

	resultingUser, err := userService.GetByID(nU.UserID)

	assert.Nil(t, err)
	assert.NotNil(t, resultingUser)
	assert.Exactly(t, nU, resultingUser)
}

func TestUser_GetByID_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	userService := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetUserByIDBase, nU.UserID)).WillReturnRows(rows)

	resultingUser, err := userService.GetByID(nU.UserID)

	assert.Nil(t, resultingUser)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestUser_GetByID_Error(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	userService := NewService(repo)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetUserByIDBase, nU.UserID)).WillReturnError(errors.New("database error"))

	resultingUser, err := userService.GetByID(nU.UserID)

	assert.Nil(t, resultingUser)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestUser_GetByEmail(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	userService := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(nU.UserID, nU.Email, nU.FirstName, nU.LastName, nU.FullName, nU.CreatedDate,
			nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.Admin, nU.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetUserByEmailBase, nU.Email)).WillReturnRows(rows)

	resultingUser, err := userService.GetByEmail(nU.Email)

	assert.Nil(t, err)
	assert.NotNil(t, resultingUser)
	assert.Exactly(t, nU, resultingUser)
}

func TestUser_GetByEmail_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	userService := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetUserByEmailBase, nU.Email)).WillReturnRows(rows)

	resultingUser, err := userService.GetByEmail(nU.Email)

	assert.Nil(t, resultingUser)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestUser_GetByEmail_Error(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	userService := NewService(repo)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetUserByEmailBase, nU.Email)).WillReturnError(errors.New("database error"))

	resultingUser, err := userService.GetByEmail(nU.Email)

	assert.Nil(t, resultingUser)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestUser_GetByAlexa(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	userService := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(nU.UserID, nU.Email, nU.FirstName, nU.LastName, nU.FullName, nU.CreatedDate,
			nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.Admin, nU.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetUserByAlexaBase, nU.AlexaUserID)).WillReturnRows(rows)

	resultingUser, err := userService.GetByAlexaID(nU.AlexaUserID)

	assert.Nil(t, err)
	assert.NotNil(t, resultingUser)
	assert.Exactly(t, nU, resultingUser)
}

func TestUser_GetByAlexa_NotFound(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	userService := NewService(repo)

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"})

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetUserByAlexaBase, nU.AlexaUserID)).WillReturnRows(rows)

	resultingUser, err := userService.GetByAlexaID(nU.AlexaUserID)

	assert.Nil(t, resultingUser)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestUser_GetByAlexa_Error(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	userService := NewService(repo)

	mock.ExpectQuery(fmt.Sprintf(dbrepo.GetUserByAlexaBase, nU.AlexaUserID)).WillReturnError(errors.New("database error"))

	resultingUser, err := userService.GetByAlexaID(nU.AlexaUserID)

	assert.Nil(t, resultingUser)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestUser_GetByAccessToken(t *testing.T) {
	db, mock, testerr := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if testerr != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, testerr)
	}
	defer db.Close()

	repo, err := dbrepo.NewRepositoryWithDB(db)
	if err != nil {
		t.Fatalf(`an error "%s" was not expected when opening the fake database connection`, err)
	}

	userService := NewService(repo)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(googleAPIOKResponse))
	})
	httpClient, teardown := testHTTPClient(h)
	defer teardown()

	client := NewClient()
	client.httpClient = httpClient

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "full_name", "created_date",
		"access_token", "refresh_token", "alexa_user_id", "is_admin", "temp_match"}).
		AddRow(nU.UserID, nU.Email, nU.FirstName, nU.LastName, nU.FullName, nU.CreatedDate,
			nU.AccessToken, nU.RefreshToken, nU.AlexaUserID, nU.Admin, nU.TempMatch)

	mock.ExpectQuery(fmt.Sprintf(`SELECT \* FROM user WHERE email = ".+"`)).WillReturnRows(rows)

	resultingUser, err := userService.GetByAccessToken(nU.AccessToken, client)

	assert.Nil(t, err)
	assert.NotNil(t, resultingUser)
}

func TestUser_Create(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestUser_GenerateTempMatch(t *testing.T) {
	assert.Equal(t, "", "")
}
