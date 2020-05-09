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

	UpdateDish(*gin.Context)
	DeleteDish(*gin.Context)

	GetStorageDishes(*gin.Context)

	GetStorageByID(*gin.Context)
	GetStorageUnits(*gin.Context)
	CreateStorageUnit(*gin.Context)
	UpdateStorageUnit(*gin.Context)
	DeleteStorageUnit(*gin.Context)

	GetUsers(*gin.Context)
	GetUserHandler(*gin.Context)
	GetUserByID(*gin.Context)
	GetUserByEmail(*gin.Context)
	CreateUser(*gin.Context)
	DeleteUser(*gin.Context)

	HandleDishesRequest(*gin.Context)
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
	RequestType string `json:"fcapiRequestType"`
	AccessToken string `json:"accessToken"`
	AlexaUserID string `json:"alexaUserID"`
	StorageID   string `json:"storageID"`
	//DishID       int    `json:"dishID"`
	Title string `json:"title"`
	//Description  string `json:"description"`
	ExpireWindow string `json:"expireWindow"`
	//Priority     string `json:"priority"`
	//DishType     string `json:"dishType"`
	//Portions     int    `json:"portions"`
}

type apiResponse struct {
	Message dishDomain.Dishes `json:"message"`
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

//------New Handler Section - Dishes-----------------------------------------------------------------------
//
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

	if aR.RequestType == "GET" {

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
			fmt.Println("NEW____-----GOT THE NORMAL GETDISHES ROUTE!!!...... in the NEW handler")
			dishID, err := strconv.Atoi(dishIDParam)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			marshaledDish, err := getDishByID(dishID, requestUser, h.dishService)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.JSON(200, gin.H{
				"message": marshaledDish,
			})
			return
		} else {
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
		}

	} else if aR.RequestType == "POST" {
		fmt.Println("doing the new createDishes() within the new dish request handler")
		err := createDish(aR, h.dishService)
		if err != nil {
			c.AbortWithStatus(err.Status())
			return
		}
		fmt.Println("Successfully added the dish to the database!")
		c.JSON(200, gin.H{
			"message": []byte("Your dish has been added to the database."),
		})
		return

	} else if aR.RequestType == "PATCH" {
		dishID, err := strconv.Atoi(dishIDParam)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		fmt.Println("got the dish update method for dish number:", dishID)

	} else if aR.RequestType == "DELETE" {

	}

	c.AbortWithStatus(http.StatusNotImplemented)
}

//---------------------------------------------------------------------------------------------------------

//Ping is the test function to see if the server is being hit.
func (h *handler) Ping(c *gin.Context) {
	fmt.Println("NEW____-----Running the Ping function: Ping")
	c.JSON(200, gin.H{
		"message": "NEW----Ping",
	})
}

func (h *handler) Pong(c *gin.Context) {
	fmt.Println("NEW - - - PONG PONG PONG - got the pong method!")
	c.JSON(200, gin.H{
		"message": "NEW----pong",
	})
}

//@@@@@@App Section@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

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
	fmt.Println("***********New interface function URL:\n", url)
	cookie := &http.Cookie{
		Name:   "oauthstate",
		Value:  oauthstate,
		MaxAge: 120,
		Secure: true,
	}
	http.SetCookie(c.Writer, cookie)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

//Oauthlogin takes a user to the external google sign in flow.
func getOAuthURL(oC oauthConfig, state string) string {
	fmt.Println("Running the getURL function which takes an oauthConfig (real or mocked), does the .AuthCodeURL() method on it")
	return oC.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

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
	fmt.Println("***********New interface function to get token:\n", token)

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

//@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

//^^^^^^^Dish Section ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

//getDishes gets all the dishes the active user has
func getDishes(requestUser *userDomain.User, service dish.Service) ([]byte, fcerr.FCErr) {
	var dishes *dishDomain.Dishes
	var err fcerr.FCErr
	fmt.Println("NEW____-----Running the GetDishes function")

	//accessToken := aR.AccessToken

	dishes, err = service.GetAll(requestUser)

	if err != nil {
		//fcerr := fcerr.NewInternalServerError("could not handle the GetDishes route")
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

func getExpiredDishes(rUser *userDomain.User, service dish.Service) ([]byte, fcerr.FCErr) {
	var dishes *dishDomain.Dishes
	var err fcerr.FCErr
	fmt.Println("NEW____-----Running the GetDishes function")

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

//CreateDish adds a dish to the list
func createDish(aR apiRequest, service dish.Service) fcerr.FCErr {

	fmt.Println("running the createDish() non-handler function")

	//1. Check if Alexa User is in DB

	//2. If so, run function with passing user id

	//3. If not, check if you can get a user in the db with the access token, and then run the function with the user id

	//4. if still not, return not authorized.

	dishMap := make(map[string]string)
	dishMap["storageID"] = aR.StorageID
	dishMap["title"] = aR.Title
	dishMap["expireWindow"] = aR.ExpireWindow

	resultingDish, err := service.Create(aR.AlexaUserID, aR.AccessToken, dishMap)

	if err != nil || resultingDish.DishID == 0 {
		return fcerr.NewInternalServerError("seems to have brokne")
	}
	return nil

}

//UpdateDish updates certain attributes of a specific dish
func (h *handler) UpdateDish(c *gin.Context) {
	dishID := c.Param("dish_id")
	fmt.Println("NEW____-----Running the UpdateDish function for this dish:", dishID)
	c.JSON(200, gin.H{
		"message": "NEW----Running the UpdateDish function for this dish:" + dishID,
	})
}

//DeleteDish deletes a specific dish from the list
func (h *handler) DeleteDish(c *gin.Context) {
	dishID := c.Param("dish_id")
	fmt.Println("NEW____-----Running the DeleteDish function for this dish:", dishID)
	c.JSON(200, gin.H{
		"message": "NEW----Running the DeleteDish function for this dish:" + dishID,
	})
}

//GetStorageDishes gets all the dishes for the active user for a specific storage unit.
func (h *handler) GetStorageDishes(c *gin.Context) {
	storageID := c.Param("storage_id")
	fmt.Println("NEW____-----Running the GetStorageDishes function for this storeage:", storageID)
	c.JSON(200, gin.H{
		"message": "NEW----Running the GetStorageDishes function for this storeage:" + storageID,
	})
}

//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

//######Storage Section#############################################################################################
func (h *handler) GetStorageByID(c *gin.Context) {
	fmt.Println("NEW____-----Running the GetStorageByID function from the new handler")
	c.JSON(200, gin.H{
		"message": "NEW----Running the GetStorageByID function from the new hanlder",
	})
}

//GetStorageUnits gets all the storage units for the active user.
func (h *handler) GetStorageUnits(c *gin.Context) {
	fmt.Println("NEW____-----Running the GetStorageUnits function")
	c.JSON(200, gin.H{
		"message": "NEW----Running the GetStorageUnits function",
	})
}

//CreateStorageUnit adds a storage unit to the list
func (h *handler) CreateStorageUnit(c *gin.Context) {
	fmt.Println("NEW____-----Running the CreateStorageUnit function")
	c.JSON(200, gin.H{
		"message": "NEW----Running the CreateStorageUnit function",
	})
}

//UpdateStorageUnit updates certain attributes of a specific storage unit
func (h *handler) UpdateStorageUnit(c *gin.Context) {
	storageID := c.Param("storage_id")
	fmt.Println("NEW____-----Running the UpdateStorageUnit function for the storage:", storageID)
	c.JSON(200, gin.H{
		"message": "NEW----Running the UpdateStorageUnit function for the storage:" + storageID,
	})
}

//DeleteStorageUnit deletes a specific storage unit from the list
func (h *handler) DeleteStorageUnit(c *gin.Context) {
	storageID := c.Param("storage_id")
	fmt.Println("NEW____-----Running the DeleteStorageUnit function for the storage:", storageID)
	c.JSON(200, gin.H{
		"message": "NEW----Running the DeleteStorageUnit function for the storage:" + storageID,
	})
}

//##################################################################################################################

//^^^^^Users Section^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

//GetUsers gets all the users the active user has permissions for.
func (h *handler) GetUsers(c *gin.Context) {
	fmt.Println("NEW____-----Running the GetUsers function")
	c.JSON(200, gin.H{
		"message": "NEW----Running the GetUsers function",
	})
}

//GetUserHandler decides if the param is an email and routes between GetUserByID and GetUserByEmail
func (h *handler) GetUserHandler(c *gin.Context) {
	userID := c.Param("dish_id")
	if userID == "expired" {
		fmt.Println("NEW____-----GOT THE GetUserByEmail ROUTE!!! ...... in the NEW handler!")
		h.GetUserByEmail(c)
	} else {
		fmt.Println("NEW____-----GOT THE NORMAL GetUserByID ROUTE!!!...... in the NEW handler")
		h.GetUserByID(c)
	}
}

//GetUserByID gets a specific user if the active user has permissions to see.
func (h *handler) GetUserByID(c *gin.Context) {
	userID := c.Param("user_id")
	fmt.Println("NEW____-----Running the GetUser function for the user with this email:", userID)
	c.JSON(200, gin.H{
		"message": "NEW----Running the GetUser function for this user:" + userID,
	})
}

//GetUserByEmail gets a specific user if the active user has permissions to see.
func (h *handler) GetUserByEmail(c *gin.Context) {
	userID := c.Param("user_id")
	fmt.Println("NEW____-----Running the GetUser function for this user:", userID)
	c.JSON(200, gin.H{
		"message": "NEW----Running the GetUser function for this user:" + userID,
	})
}

//CreateUser adds a user to the list
func (h *handler) CreateUser(c *gin.Context) {
	fmt.Println("NEW____-----Running the CreateUser function")
	c.JSON(200, gin.H{
		"message": "NEW----Running the CreateUser function",
	})
}

//DeleteUser deletes a specific user from the list
func (h *handler) DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	fmt.Println("NEW____-----Running the DeleteUser function for this user:", userID)
	c.JSON(200, gin.H{
		"message": "NEW----Running the DeleteUser function for this user:" + userID,
	})
}

//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

func numGenerator() string {
	n := make([]byte, 32)
	rand.Read(n)
	fmt.Println("Old way:", base64.StdEncoding.EncodeToString(n))
	fmt.Println("New way:", base64.URLEncoding.EncodeToString(n))

	return base64.URLEncoding.EncodeToString(n)

}
