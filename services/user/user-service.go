package user

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/jasonradcliffe/freshness-countdown-api/domain/user"
	"github.com/jasonradcliffe/freshness-countdown-api/repository/db"

	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
)

//Service is the interface that defines the contract for a dish service.
type Service interface {
	GetByID(int) (*user.User, fcerr.FCErr)
	GetByEmail(string) (*user.User, fcerr.FCErr)
	GetByAlexaID(string) (*user.User, fcerr.FCErr)
	GetByAccessToken(string) (*user.User, fcerr.FCErr)
	Create(u user.OauthUser, aT string, rT string) (*user.User, fcerr.FCErr)
}

type service struct {
	repository db.Repository
}

//NewService takes a database repository and gives you a new Service instance.
func NewService(repo db.Repository) Service {
	return &service{
		repository: repo,
	}
}

//GetByID gets a user from the database with the given ID
func (s *service) GetByID(id int) (*user.User, fcerr.FCErr) {
	return nil, nil
}

//GetByEmail gets a user from the database with the given email address
func (s *service) GetByEmail(email string) (*user.User, fcerr.FCErr) {
	receivedUser, err := s.repository.GetUserByEmail(email)
	if err != nil {
		fmt.Println("user service could not get the user by email")
		fcerr := fcerr.NewInternalServerError("user service could not get the user by email")
		return nil, fcerr
	}

	return receivedUser, nil
}

//GetByAlexaID gets a user from the database with the given alexa user id
func (s *service) GetByAlexaID(alexaID string) (*user.User, fcerr.FCErr) {
	receivedUser, err := s.repository.GetUserByEmail(alexaID)
	if err != nil {
		fmt.Println("user service could not get the user by email")
		fcerr := fcerr.NewInternalServerError("user service could not get the user by email")
		return nil, fcerr
	}

	return receivedUser, nil
}

//GetByAccessToken gets a user from the database with the given access token
func (s *service) GetByAccessToken(aT string) (*user.User, fcerr.FCErr) {
	receivedUser, err := s.repository.GetUserByEmail(aT)
	if err != nil {
		fmt.Println("user service could not get the user by email")
		fcerr := fcerr.NewInternalServerError("user service could not get the user by email")
		return nil, fcerr
	}

	return receivedUser, nil
}

func (s *service) Create(u user.OauthUser, aT string, rT string) (*user.User, fcerr.FCErr) {
	var newUser user.User
	createdDate := "2016-01-02T15:04:05"
	tempMatch := s.GenerateTempMatch()

	newUser.Email = u.Email
	newUser.FirstName = u.FirstName
	newUser.LastName = u.LastName
	newUser.FullName = u.FullName
	newUser.CreatedDate = createdDate
	newUser.AccessToken = aT
	newUser.RefreshToken = rT
	newUser.TempMatch = tempMatch

	receivedUser, err := s.repository.CreateUser(newUser)
	if err != nil {
		fmt.Println("the user service could not create the new user")
		fcerr := fcerr.NewInternalServerError("the user service could not create the new user")
		return nil, fcerr
	}

	return receivedUser, nil
}

func (s *service) GenerateTempMatch() string {
	n := make([]byte, 15)
	rand.Read(n)
	return base64.URLEncoding.EncodeToString(n)
}
