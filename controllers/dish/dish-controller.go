package dish

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//GetHandler differentiates between getDishes and getExpired as Httprouter is not able to.
func GetHandler(c *gin.Context) {
	fmt.Println("Running the GetHandler function: beep boop")
	dishID := c.Param("dish_id")
	if dishID == "expired" {
		fmt.Println("GOT THE EXPIRED ROUTE!!!")
		GetExpiredDishes(c)
	} else {
		fmt.Println("GOT THE NORMAL GETDISHES ROUTE!!!")
		GetDish(c)
	}
}

//GetDishes gets all the dishes for the active user.
func GetDishes(c *gin.Context) {
	fmt.Println("Running the GetDishes function")
	c.JSON(200, gin.H{
		"message": "Running the GetDishes function",
	})
}

//GetExpiredDishes gets all the dishes for the active user that are expired.
func GetExpiredDishes(c *gin.Context) {
	fmt.Println("Running the GetExpiredDishes function")
	c.JSON(200, gin.H{
		"message": "Running the GetExpiredDishes function",
	})
}

//GetDish gets a specific dish.
func GetDish(c *gin.Context) {
	dishID := c.Param("dish_id")
	fmt.Println("Running the GetDish function for this dish:", dishID)
	c.JSON(200, gin.H{
		"message": "Running the GetDish function for this dish:" + dishID,
	})
}

//CreateDish adds a dish to the list
func CreateDish(c *gin.Context) {
	fmt.Println("Running the CreateDish function")
	c.JSON(200, gin.H{
		"message": "Running the CreateDish function",
	})
}

//UpdateDish updates certain attributes of a specific dish
func UpdateDish(c *gin.Context) {
	dishID := c.Param("dish_id")
	fmt.Println("Running the UpdateDish function for this dish:", dishID)
	c.JSON(200, gin.H{
		"message": "Running the UpdateDish function for this dish:" + dishID,
	})
}

//DeleteDish deletes a specific dish from the list
func DeleteDish(c *gin.Context) {
	dishID := c.Param("dish_id")
	fmt.Println("Running the DeleteDish function for this dish:", dishID)
	c.JSON(200, gin.H{
		"message": "Running the DeleteDish function for this dish:" + dishID,
	})
}

//GetStorageDishes gets all the dishes for the active user for a specific storage unit.
func GetStorageDishes(c *gin.Context) {
	storageID := c.Param("storage_id")
	fmt.Println("Running the GetStorageDishes function for this storeage:", storageID)
	c.JSON(200, gin.H{
		"message": "Running the GetStorageDishes function for this storeage:" + storageID,
	})
}
