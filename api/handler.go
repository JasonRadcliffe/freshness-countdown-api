package api

import (
	"fmt"

	"github.com/jasonradcliffe/freshness-countdown-api/fcerr"

	"github.com/gin-gonic/gin"
	dishDomain "github.com/jasonradcliffe/freshness-countdown-api/domain/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/services/dish"

	"github.com/jasonradcliffe/freshness-countdown-api/services/storage"
)

//Handler interface is the contract for the methods that the handler needs to have.
type Handler interface {
	Ping(*gin.Context)

	GetDishes(*gin.Context)
	GetDishHandler(*gin.Context)
	GetDishByID(*gin.Context)
	GetExpiredDishes(*gin.Context)
	CreateDish(*gin.Context)
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
	UpdateUser(*gin.Context)
	DeleteUser(*gin.Context)
}

type handler struct {
	dishService    dish.Service
	storageService storage.Service
}

//NewHandler takes a sequence of services and returns a new API Handler.
func NewHandler(ds dish.Service, ss storage.Service) Handler {
	return &handler{
		dishService:    ds,
		storageService: ss,
	}
}

//Ping is the test function to see if the server is being hit.
func (h *handler) Ping(c *gin.Context) {
	fmt.Println("NEW____-----Running the Ping function: PONG")
	c.JSON(200, gin.H{
		"message": "NEW----pong",
	})
}

//^^^^^^^Dish Section ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

//GetDishes gets all the dishes the active user has
func (h *handler) GetDishes(c *gin.Context) {
	var dishes *dishDomain.Dishes
	var err fcerr.FCErr
	fmt.Println("NEW____-----Running the GetDishes function")

	dishes, err = h.dishService.GetAll()
	if err != nil {
		//fcerr := fcerr.NewInternalServerError("could not handle the GetDishes route")
		fmt.Println("could not handle the GetDishes route")
		return
	}
	fmt.Println("I think we got some dishes!!! The first of which is:", (*dishes)[0])
	c.JSON(200, gin.H{
		"message": "NEW----Running the GetDishes function",
	})
}

func (h *handler) GetDishHandler(c *gin.Context) {
	dishID := c.Param("dish_id")
	if dishID == "expired" {
		fmt.Println("NEW____-----GOT THE EXPIRED ROUTE!!! ...... in the NEW handler!")
		h.GetExpiredDishes(c)
	} else {
		fmt.Println("NEW____-----GOT THE NORMAL GETDISHES ROUTE!!!...... in the NEW handler")
		h.GetDishByID(c)
	}
}

//GetDishByID gets a specific dish if it belongs to the current user
func (h *handler) GetDishByID(c *gin.Context) {
	dishID := c.Param("dish_id")
	fmt.Println("NEW____-----Running the GetDishByID function from the new handler for this dish:", dishID)
	c.JSON(200, gin.H{
		"message": "NEW----Running the GetDish function from the new handler for this dish:" + dishID,
	})
}

//GetExpiredDishes gets all the dishes for the current user that are expired
func (h *handler) GetExpiredDishes(c *gin.Context) {
	fmt.Println("NEW____-----Running the GetExpiredDishes function from the new handler")
	c.JSON(200, gin.H{
		"message": "NEW----Running the GetExpiredDishes function from the new hanlder",
	})
}

//CreateDish adds a dish to the list
func (h *handler) CreateDish(c *gin.Context) {
	fmt.Println("NEW____-----Running the CreateDish function")
	c.JSON(200, gin.H{
		"message": "NEW----Running the CreateDish function",
	})
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
