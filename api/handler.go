package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"
	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	dishDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	storageDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/storage"
	userDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/user"
	"github.com/jasonradcliffe/freshness-countdown-api/services/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/services/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/services/user"
)

//Handler interface is the contract for the methods that the handler needs to have.
type Handler interface {
	Ping(*gin.Context)
	Pong(*gin.Context)

	Login(*gin.Context)
	Oauthlogin(*gin.Context)
	LoginSuccess(*gin.Context)

	HandleDishesRequest(*gin.Context)
	DishExpiresIn(*gin.Context)
	DishExpiresBy(*gin.Context)

	HandleStorageRequest(*gin.Context)
	HandleUsersRequest(*gin.Context)
}

type oauthConfig interface {
	//Exchange func(ctx context.Context, code string, opts ...AuthCodeOption)
	Exchange(context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	//AuthCodeURL func(state string, opts ...AuthCodeOption)
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
}

type handler struct {
	dishService    dish.Service
	storageService storage.Service
	userService    user.Service
	oauthConfig    oauthConfig
}

type apiRequest struct {
	RequestType  string `json:"fcapiRequestType"`
	AccessToken  string `json:"accessToken"`
	AlexaUserID  string `json:"alexaUserID"`
	StorageID    string `json:"storageID"`
	DishID       int    `json:"dishID"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ExpireWindow string `json:"expireWindow"`
	Priority     string `json:"priority"`
	DishType     string `json:"dishType"`
	Portions     int    `json:"portions"`
}

var oauthstate string
var currentUser userDomain.OauthUser

//NewHandler takes a sequence of services and returns a new API Handler.
func NewHandler(ds dish.Service, ss storage.Service, us user.Service, oC oauthConfig) Handler {
	return &handler{
		dishService:    ds,
		storageService: ss,
		userService:    us,
		oauthConfig:    oC,
	}
}

//ValidateUser looks at the request details and extracts the user making the request. Err is returned if not able to find OR add a user
func ValidateUser(h *handler, aR apiRequest) (*userDomain.User, fcerr.FCErr) {
	alexaIDUser, err := h.userService.GetByAlexaID(aR.AlexaUserID)
	if err != nil {
		fmt.Println("couldn't get a user from alexa id:" + aR.AlexaUserID)
		accessTokenUser, err2 := h.userService.GetOrCreateByAccessToken(aR.AccessToken, user.NewClient())
		if err2 != nil {
			fmt.Println("couldn't get or create a user with access token:" + aR.AccessToken)
			return nil, fcerr.NewUnauthorizedError("Could not validate this user")
		}

		fmt.Println("Here is the user we got from the access token!" + accessTokenUser.Email)
		fmt.Println("We should add the user's alexa ID since we know the db doesn't have it")
		_, err3 := h.userService.UpdateAlexaID(*accessTokenUser, aR.AlexaUserID)
		if err3 != nil {
			fmt.Println("We couldn't add the alexa user id of the new user - no biggie")
		}
		return accessTokenUser, nil
	}

	fmt.Println("Here is the user we got from the Alexa ID!" + alexaIDUser.Email)
	return alexaIDUser, nil
}

//------Dishes Handler and Helpers---------------------------------------------------------------------------------------------------------------
func (h *handler) HandleDishesRequest(c *gin.Context) {
	var aR apiRequest

	if err := c.ShouldBindJSON(&aR); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if aR.AlexaUserID == "" && aR.AccessToken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	requestUser, err := ValidateUser(h, aR)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	dishIDParam := c.Param("dish_id")

	switch aR.RequestType {

	case "GET":
		if dishIDParam == "expired" {
			fmt.Println("got the post request for GetExpired!")
			marshaledDishList, err := getExpiredDishes(requestUser, h.dishService)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.JSON(200, gin.H{
				"message": marshaledDishList,
			})
			return
		} else if dishIDParam != "" {
			fmt.Println("hitting the getDish handler")
			dishID, err := strconv.Atoi(dishIDParam)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			fmt.Println("dishID:" + strconv.Itoa(dishID))

			marshaledDish, err := getDishByID(dishID, requestUser, h.dishService)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.JSON(200, gin.H{
				"message": marshaledDish,
			})
			return
		}

		fmt.Println("got the getDishes route!!!")
		marshaledDishList, err := getDishes(requestUser, h.dishService)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(200, gin.H{
			"message": marshaledDishList,
		})
		return

	case "POST":
		fmt.Println("doing the createDishes() within the dish request handler")
		err := createDish(requestUser, aR, h.dishService)
		if err != nil {
			c.AbortWithStatus(err.Status())
			return
		}
		fmt.Println("Successfully added the dish to the database!")
		c.JSON(200, gin.H{
			"message": []byte("Your dish has been added to the database."),
		})
		return

	case "PATCH":
		dishID, err := strconv.Atoi(dishIDParam)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		fmt.Println("got the dish update method for dish number:", dishID)
		err2 := updateDish(requestUser, aR, h.dishService)
		if err2 != nil {
			fmt.Println("Got an error when doing the update dish route")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		fmt.Println("Successfully updated the dish in the database!")
		c.JSON(200, gin.H{
			"message": []byte("Your dish has been updated in the database."),
		})
	case "DELETE":
		dishID, err := strconv.Atoi(dishIDParam)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		fmt.Println("got the dish delete method for dish number:", dishID)
		err2 := deleteDish(requestUser, dishID, h.dishService)
		if err2 != nil {
			fmt.Println("Got an error when doing the update dish route")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		fmt.Println("Successfully deleted the dish from the database!")
		c.JSON(200, gin.H{
			"message": []byte("Your dish has been deleted from the database."),
		})

	default:
		c.AbortWithStatus(http.StatusNotImplemented)

	}
}

func (h *handler) DishExpiresIn(c *gin.Context) {
	expirationWindow := c.Param("duration")

	if expirationWindow == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fmt.Println("got the DishExpiresIn handler")
	c.JSON(200, gin.H{
		"message": []byte("not yet implemented the DishExpiresIn handler"),
	})

}

func (h *handler) DishExpiresBy(c *gin.Context) {
	date := c.Param("date")

	if date == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fmt.Println("got the DishExpiresBy handler")
	c.JSON(200, gin.H{
		"message": []byte("not yet implemented the DishExpiresBy handler"),
	})

}

//getDishes gets all the dishes the active user has
func getDishes(requestUser *userDomain.User, service dish.Service) ([]byte, fcerr.FCErr) {
	var dishes *dishDomain.Dishes
	var err fcerr.FCErr
	fmt.Println("Running the getDishes function")

	//accessToken := aR.AccessToken

	dishes, err = service.GetAll(requestUser)

	if err != nil {
		fmt.Println("could not handle the GetDishes route")
		return nil, fcerr.NewInternalServerError("unsuccessful at service.GetAll")
	}

	fmt.Println("I think we got some dishes!!! The first of which is:", (*dishes)[0])

	marshaledDishes, merr := json.Marshal(dishes)
	if merr != nil {
		return nil, fcerr.NewInternalServerError("JSON Error - Could not marshal the dishes")
	}
	return marshaledDishes, nil
}

//getDishByID gets a particular dish the requesting user has
func getDishByID(pID int, requestingUser *userDomain.User, service dish.Service) ([]byte, fcerr.FCErr) {
	var dish *dishDomain.Dish
	var err fcerr.FCErr
	fmt.Println("running non-gin getDishByID func")

	//accessToken := aR.AccessToken

	dish, err = service.GetByID(requestingUser, pID)

	if err != nil {
		//fcerr := fcerr.NewInternalServerError("could not handle the GetDishes route")
		fmt.Println("could not handle the GetDishes route")
		return nil, fcerr.NewInternalServerError("unsuccessful at service.GetAll")
	}

	fmt.Println("I think we got a dish!!! It is:", dish.Title)

	marshaledDish, merr := json.Marshal(dish)
	if merr != nil {
		return nil, fcerr.NewInternalServerError("JSON Error - Could not marshal the dishes")
	}
	return marshaledDish, nil
}

//getExpiredDishes gets all the dishes the active user has that have already expired
func getExpiredDishes(rUser *userDomain.User, service dish.Service) ([]byte, fcerr.FCErr) {
	var dishes *dishDomain.Dishes
	var err fcerr.FCErr
	fmt.Println("Running the getExpiredDishes function")

	//accessToken := aR.AccessToken

	dishes, err = service.GetExpired(rUser)

	if err != nil {
		//fcerr := fcerr.NewInternalServerError("could not handle the GetDishes route")
		fmt.Println("could not handle the get expired dishes handle function")
		return nil, fcerr.NewInternalServerError("unsuccessful at service.GetAll")
	}

	fmt.Println("I think we got some dishes!!! The first of which is:", (*dishes)[0])
	fmt.Println("The length of the list we got is:", len(*dishes))

	marshaledDishes, merr := json.Marshal(dishes)
	if merr != nil {
		return nil, fcerr.NewInternalServerError("JSON Error - Could not marshal the dishes")
	}
	return marshaledDishes, nil
}

//createDish adds a dish to the list
func createDish(requestingUser *userDomain.User, aR apiRequest, service dish.Service) fcerr.FCErr {

	fmt.Println("running the createDish() non-handler function")

	storageID, err := strconv.Atoi(aR.StorageID)
	if err != nil {
		return fcerr.NewBadRequestError("Error when creating the dish.")
	}

	newDish := &dishDomain.Dish{
		StorageID:   storageID,
		Title:       aR.Title,
		Description: aR.Description,
		Priority:    aR.Priority,
		DishType:    aR.DishType,
		Portions:    aR.Portions,
	}
	expireWindow := aR.ExpireWindow

	resultingDish, err := service.Create(requestingUser, newDish, expireWindow)

	if err != nil || resultingDish.DishID == 0 {
		return fcerr.NewInternalServerError("seems to have brokne")
	}
	return nil

}

//updateDish takes a requesting user, and an API request along with the dish service to update the dish to the values contained in the apirequest
func updateDish(requestingUser *userDomain.User, aR apiRequest, service dish.Service) fcerr.FCErr {
	fmt.Println("running the updateDish() non-handler function")

	storageID, err := strconv.Atoi(aR.StorageID)
	if err != nil {
		return fcerr.NewBadRequestError("Error when creating the dish.")
	}

	newDish := &dishDomain.Dish{
		StorageID:   storageID,
		Title:       aR.Title,
		Description: aR.Description,
		Priority:    aR.Priority,
		DishType:    aR.DishType,
		Portions:    aR.Portions,
	}
	expireWindow := aR.ExpireWindow

	err2 := service.Update(requestingUser, newDish, expireWindow)

	if err2 != nil {
		return fcerr.NewInternalServerError("Error when updating the dish")
	}
	return nil
}

//deleteDish takes a requesting user, and a dish ID along with the dish service to delete the dish with the personal id given
func deleteDish(requestingUser *userDomain.User, dishID int, service dish.Service) fcerr.FCErr {
	fmt.Println("running the updateDish() non-handler function")
	err := service.Delete(requestingUser, dishID)
	if err != nil {
		return fcerr.NewInternalServerError("Error when deleting the dish")
	}
	return nil
}

//---------------------------------------------------------------------------------------------------------------------------------------------------

//****Storage Handler and helpers********************************************************************************************************************
func (h *handler) HandleStorageRequest(c *gin.Context) {
	var aR apiRequest

	if err := c.ShouldBindJSON(&aR); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if aR.AlexaUserID == "" && aR.AccessToken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	requestUser, err := ValidateUser(h, aR)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	storageIDParam := c.Param("storage_id")

	switch aR.RequestType {

	case "GET":
		if storageIDParam != "" {
			fmt.Println("GOT THE NORMAL GETStorage ROUTE!!!")
			storageID, err := strconv.Atoi(storageIDParam)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			marshaledStorage, err := getStorageByID(storageID, requestUser, h.storageService)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.JSON(200, gin.H{
				"message": marshaledStorage,
			})
			return
		}
		fmt.Println("got the getStorage route!!!")
		marshaledStorageList, err := getStorage(requestUser, h.storageService)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(200, gin.H{
			"message": marshaledStorageList,
		})
		return

	case "POST":
		fmt.Println("doing the createStorage() within the storage request handler")
		err := createStorage(requestUser, aR, h.storageService)
		if err != nil {
			c.AbortWithStatus(err.Status())
			return
		}
		fmt.Println("Successfully added the storage to the database!")
		c.JSON(200, gin.H{
			"message": []byte("Your storage unit has been added to the database."),
		})
		return

	case "PATCH":
		storageID, err := strconv.Atoi(storageIDParam)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		fmt.Println("got the storage update method for storage number:", storageID)
		err2 := updateStorage(requestUser, aR, h.storageService)
		if err2 != nil {
			fmt.Println("Got an error when doing the update storage route")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	case "DELETE":
		storageID, err := strconv.Atoi(storageIDParam)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		fmt.Println("got the storage delete method for storage number:", storageID)
		err2 := deleteStorage(requestUser, storageID, h.storageService)
		if err2 != nil {
			fmt.Println("Got an error when doing the update storage route")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	default:
		c.AbortWithStatus(http.StatusNotImplemented)

	}

}

//getStorage gets all the storage units the requesting user has
func getStorage(requestUser *userDomain.User, service storage.Service) ([]byte, fcerr.FCErr) {
	var storageList *storageDomain.Storages
	var err fcerr.FCErr
	fmt.Println("Running the getStorage function")

	//accessToken := aR.AccessToken

	storageList, err = service.GetAll(requestUser)

	if err != nil {
		//fcerr := fcerr.NewInternalServerError("could not handle the GetStorage route")
		fmt.Println("could not handle the GetStorage route")
		return nil, fcerr.NewInternalServerError("unsuccessful at service.GetAll")
	}

	fmt.Println("I think we got some storage units!!! The first of which is:", (*storageList)[0])

	marshaledStorageList, merr := json.Marshal(storageList)
	if merr != nil {
		return nil, fcerr.NewInternalServerError("JSON Error - Could not marshal the storage units")
	}
	return marshaledStorageList, nil
}

//getStorageByID gets a particular storage unit the requesting user has
func getStorageByID(pID int, requestingUser *userDomain.User, service storage.Service) ([]byte, fcerr.FCErr) {
	var storage *storageDomain.Storage
	var err fcerr.FCErr
	fmt.Println("running non-gin getStorageByID func")

	storage, err = service.GetByID(requestingUser, pID)

	if err != nil {
		//fcerr := fcerr.NewInternalServerError("could not handle the GetStorageByID route")
		fmt.Println("could not handle the GetStorageByID route")
		return nil, fcerr.NewInternalServerError("unsuccessful at service.GetAll")
	}

	fmt.Println("I think we got a storage unit!!! It is:", storage.Title)

	marshaledStorage, merr := json.Marshal(storage)
	if merr != nil {
		return nil, fcerr.NewInternalServerError("JSON Error - Could not marshal the storage units")
	}
	return marshaledStorage, nil
}

//createStorage adds a storage unit to the list
func createStorage(requestingUser *userDomain.User, aR apiRequest, service storage.Service) fcerr.FCErr {

	fmt.Println("running the createStorage() non-handler function")

	newStorage := &storageDomain.Storage{
		Title:       aR.Title,
		Description: aR.Description,
	}

	resultingStorage, err := service.Create(requestingUser, newStorage)

	if err != nil || resultingStorage.StorageID == 0 {
		return fcerr.NewInternalServerError("seems to have brokne")
	}
	return nil

}

//updateStorage takes a requesting user, and an API request along with the storage service to update the storage unit to the values contained in the apirequest
func updateStorage(requestingUser *userDomain.User, aR apiRequest, service storage.Service) fcerr.FCErr {
	fmt.Println("running the updateStorage() function")

	storageID, err := strconv.Atoi(aR.StorageID)
	if err != nil {
		return fcerr.NewBadRequestError("Error when creating the storage unit.")
	}

	newStorage := &storageDomain.Storage{
		StorageID:   storageID,
		Title:       aR.Title,
		Description: aR.Description,
	}

	err2 := service.Update(requestingUser, newStorage)

	if err2 != nil {
		return fcerr.NewInternalServerError("Error when updating the storage unit")
	}
	return nil
}

//deleteStorage takes a requesting user, and a storage ID along with the storage service to delete the storage unit with the personal id given
func deleteStorage(requestingUser *userDomain.User, storageID int, service storage.Service) fcerr.FCErr {
	fmt.Println("running the deleteStorage() function")
	err := service.Delete(requestingUser, storageID)
	if err != nil {
		return fcerr.NewInternalServerError("Error when deleting the storage unit")
	}
	return nil
}

//*****************************************************************************************************************************************************

//^^^^^^^^^Users Handler and helpers^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
func (h *handler) HandleUsersRequest(c *gin.Context) {
	var aR apiRequest

	if err := c.ShouldBindJSON(&aR); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if aR.AlexaUserID == "" && aR.AccessToken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	requestUser, err := ValidateUser(h, aR)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	switch aR.RequestType {

	case "PATCH":
		fmt.Println("doing the updateUsers() within the users request handler for this user:", requestUser.Email)
		c.JSON(200, gin.H{
			"message": []byte("Your user has been updated in the database."),
		})
		return
	case "DELETE":
		fmt.Println("doing the deleteUsers() within the users request handler for this user:", requestUser.Email)
		c.JSON(200, gin.H{
			"message": []byte("Your user has been removed from the database."),
		})
		return

	default:
		c.AbortWithStatus(http.StatusNotImplemented)

	}
}

//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

//Ping is the test function to see if the server is being hit.
func (h *handler) Ping(c *gin.Context) {
	fmt.Println("Running the Ping function: Ping")
	c.JSON(200, gin.H{
		"message": "ping says: PONG",
	})
}

func (h *handler) Pong(c *gin.Context) {
	fmt.Println("got the pong method!")
	c.JSON(200, gin.H{
		"message": "pong says: PING",
	})
}

//@@@@@@App Handlers@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

//Login displays a simple link that takes a user to the external google sign in flow.
func (h *handler) Login(c *gin.Context) {
	fmt.Println("Running the Login function")
	siteData := []byte("<a href=/oauthlogin> Login with Google </a>")
	c.Data(200, "text/html", siteData)
}

//Oauthlogin takes a user to the external google sign in flow.
func (h *handler) Oauthlogin(c *gin.Context) {
	fmt.Println("Running the Oauthlogin function")
	oauthstate := numGenerator()
	url := getOAuthURL(h.oauthConfig, oauthstate)
	cookie := &http.Cookie{
		Name:   "oauthstate",
		Value:  oauthstate,
		MaxAge: 120,
		Secure: true,
	}
	http.SetCookie(c.Writer, cookie)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

//getOAuthURL takes an oauthConfig (real or mocked) and does the .AuthCodeURL() with it.
func getOAuthURL(oC oauthConfig, state string) string {
	return oC.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

//getOAuthToken takes an oauthConfig (real or mocked) and does the .Exchange() with it.
func getOAuthToken(oC oauthConfig, c *gin.Context, code string) (*oauth2.Token, error) {
	token, err := oC.Exchange(c, code)
	return token, err

}

//LoginSuccess is where the Oauth provider routes to after successfully authenticating a user
func (h *handler) LoginSuccess(c *gin.Context) {

	receivedCookie, err := c.Cookie("oauthstate")
	if err != nil {
		fmt.Println("got an error when retrieving the cookie during loginSuccess()")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println("In LoginSuccess - got the cookie:", receivedCookie)

	receivedState := c.Request.FormValue("state")
	if receivedState != receivedCookie {
		fmt.Println("receivedState:", receivedState, "did not equal oauthstate:", oauthstate)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	code := c.Request.FormValue("code")
	token, err := getOAuthToken(h.oauthConfig, c, code)
	if err != nil {
		fmt.Println("error when exchanging the token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response, err := http.Get("https://openidconnect.googleapis.com/v1/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("error when getting the userinfo with the access token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	json.Unmarshal(contents, &currentUser)
	fmt.Println("Here is the current User:", currentUser)

	if currentUser.VerifiedEmail == false {
		fmt.Println("current user.VerifiedEmail is false. CurrentUser:", currentUser)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	fmt.Println("Got a verified user!!!!!!", currentUser)

	dbUser, err := h.userService.GetByEmail(currentUser.Email)
	if err != nil {
		fmt.Println("was not able to check the database for the user on login success")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	} else if dbUser.UserID <= 0 {
		fmt.Println("loginSuccess could not find this user in the database! We should add them!!")
		receivedUser, err := h.userService.Create(currentUser, token.AccessToken, token.RefreshToken)
		if err != nil {
			fmt.Println("Was not successful in adding a new user to the database!")
			c.AbortWithStatus(http.StatusInternalServerError)
			return

		}
		fmt.Println("we just put a new user in the database!! with database user id:", receivedUser.UserID)

	}
	fmt.Println("We already have this user!!! database user id:", dbUser)

	successData := []byte("<h1>Success!</h1>")
	c.Data(200, "text/html", successData)

}

//@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

func numGenerator() string {
	n := make([]byte, 32)
	rand.Read(n)
	fmt.Println("Old way:", base64.StdEncoding.EncodeToString(n))
	fmt.Println("New way:", base64.URLEncoding.EncodeToString(n))

	return base64.URLEncoding.EncodeToString(n)

}
