package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

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
	UpdateAlexaID(user.User, string) (*user.User, fcerr.FCErr)
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
	receivedUser, err := s.repository.GetUserByID(id)
	if err != nil && err.Status() == http.StatusNotFound {
		return nil, fcerr.NewNotFoundError("Could not find this user in the system.")
	} else if err != nil {
		return nil, fcerr.NewInternalServerError("Error while retrieving the user.")
	}

	return receivedUser, nil
}

//GetByEmail gets a user from the database with the given email address
func (s *service) GetByEmail(email string) (*user.User, fcerr.FCErr) {
	receivedUser, err := s.repository.GetUserByEmail(email)
	if err != nil && err.Status() == http.StatusNotFound {
		fmt.Println("user service could not get the user by email")
		fcerr := fcerr.NewNotFoundError("Did not find this user.")
		return nil, fcerr
	} else if err != nil {
		fcerr := fcerr.NewInternalServerError("user service could not get the user by email")
		return nil, fcerr
	}

	return receivedUser, nil
}

//GetByAlexaID gets a user from the database with the given alexa user id
func (s *service) GetByAlexaID(alexaID string) (*user.User, fcerr.FCErr) {
	receivedUser, err := s.repository.GetUserByAlexa(alexaID)
	if err != nil {
		fmt.Println("user service could not get the user by alexa ID")
		fcerr := fcerr.NewNotFoundError("user service could not get the user by email")
		return nil, fcerr
	}

	return receivedUser, nil
}

//GetByAccessToken gets a user from the database with the given access token
func (s *service) GetByAccessToken(aT string) (*user.User, fcerr.FCErr) {

	var currentUser user.OauthUser

	response, err := http.Get("https://openidconnect.googleapis.com/v1/userinfo?access_token=" + aT)
	if err != nil {
		fmt.Println("error when getting the userinfo with the access token")
		return nil, fcerr.NewInternalServerError("Error when trying to verify user identity")
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fcerr.NewInternalServerError("Error when trying to read response from Google about user identity")
	}
	fmt.Println("\nHere is the contents:\n", contents)

	json.Unmarshal(contents, &currentUser)
	fmt.Println("Here is the current User we are fetching with access token:", currentUser)

	if currentUser.VerifiedEmail == false {
		fmt.Println("current user.VerifiedEmail is false. CurrentUser:", currentUser)
		return nil, fcerr.NewBadRequestError("Not Authorized. Please verify email address.")
	}

	fmt.Println("Got a verified user!!!!!!", currentUser)

	dbUser, err2 := s.GetByEmail(currentUser.Email)
	if err2 != nil && err2.Status() == http.StatusNotFound {
		fmt.Println("We could not find this user in the database! (We should add them!?!)")
		newUser, err := s.Create(currentUser, aT, "")
		if err != nil {
			return nil, fcerr.NewInternalServerError("Attempted to add the user to the database, but something went wrong.")
		}
		fmt.Println("User has been added. New User ID:" + strconv.Itoa(newUser.UserID))
		return newUser, nil

	} else if err2 != nil || dbUser.UserID <= 0 {
		fmt.Println("was not able to check the database for the user on login success")
		return nil, fcerr.NewInternalServerError("Was not able to check for the user after getting email address.")
	}

	fmt.Println("We already have this user!!! database user id:", dbUser)
	return dbUser, nil

}

func (s *service) Create(u user.OauthUser, aT string, rT string) (*user.User, fcerr.FCErr) {
	var newUser user.User

	timeNow := time.Now().In(time.UTC)
	createdDate := timeNow.Format("2006-01-02T15:04:05")

	newUser.Email = u.Email
	newUser.FirstName = u.FirstName
	newUser.LastName = u.LastName
	newUser.FullName = u.FullName
	newUser.CreatedDate = createdDate
	newUser.AccessToken = aT
	newUser.RefreshToken = rT

	receivedUser, err := s.repository.CreateUser(newUser)
	if err != nil {
		fmt.Println("the user service could not create the new user")
		fcerr := fcerr.NewInternalServerError("the user service could not create the new user")
		return nil, fcerr
	}

	return receivedUser, nil
}

//UpdateAlexaID (u user.User, alexaID string) takes a user and a string and sets the alexaUserID equal to the given string in the database
func (s *service) UpdateAlexaID(u user.User, alexaID string) (*user.User, fcerr.FCErr) {
	newUser := &user.User{
		UserID:       u.UserID,
		Email:        u.Email,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		FullName:     u.FullName,
		CreatedDate:  u.CreatedDate,
		AccessToken:  u.AccessToken,
		RefreshToken: u.RefreshToken,
		AlexaUserID:  alexaID,
	}
	updatedUser, err := s.repository.UpdateUser(*newUser)
	if err != nil {
		return nil, fcerr.NewInternalServerError("Error when updating the user with Alexa ID")
	}

	return updatedUser, nil
}
