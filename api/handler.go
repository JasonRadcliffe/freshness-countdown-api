package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jasonradcliffe/freshness-countdown-api/services/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/services/storage"
)

//Handler interface is the contract for the methods that the handler needs to have.
type Handler interface {
	GetDishHandler(*gin.Context)
	GetDishByID(*gin.Context)
	GetStorageByID(*gin.Context)
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

func (h *handler) GetDishHandler(c *gin.Context) {
	dishID := c.Param("dish_id")
	if dishID == "expired" {
		fmt.Println("GOT THE EXPIRED ROUTE!!! ...... in the NEW handler!")
		h.GetExpiredDishes(c)
	} else {
		fmt.Println("GOT THE NORMAL GETDISHES ROUTE!!!...... in the NEW handler")
		h.GetDishByID(c)
	}
}

func (h *handler) GetDishByID(c *gin.Context) {
	dishID := c.Param("dish_id")
	fmt.Println("Running the GetDishByID function from the new handler for this dish:", dishID)
	c.JSON(200, gin.H{
		"message": "Running the GetDish function from the new handler for this dish:" + dishID,
	})
}

func (h *handler) GetExpiredDishes(c *gin.Context) {
	fmt.Println("Running the GetExpiredDishes function from the new handler")
	c.JSON(200, gin.H{
		"message": "Running the GetExpiredDishes function from the new hanlder",
	})
}

func (h *handler) GetStorageByID(c *gin.Context) {
	fmt.Println("Running the GetStorageByID function from the new handler")
	c.JSON(200, gin.H{
		"message": "Running the GetStorageByID function from the new hanlder",
	})
}
