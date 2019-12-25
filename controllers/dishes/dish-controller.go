package dishes

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//GetDishes gets all the dishes for the active user.
func GetDishes(c *gin.Context) {
	fmt.Println("Running the GetDishes function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//GetExpiredDishes gets all the dishes for the active user that are expired.
func GetExpiredDishes(c *gin.Context) {
	fmt.Println("Running the GetExpiredDishes function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//GetDish gets a specific dish.
func GetDish(c *gin.Context) {
	fmt.Println("Running the GetDish function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//CreateDish adds a dish to the list
func CreateDish(c *gin.Context) {
	fmt.Println("Running the CreateDish function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//UpdateDish updates certain attributes of a specific dish
func UpdateDish(c *gin.Context) {
	fmt.Println("Running the UpdateDish function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//DeleteDish deletes a specific dish from the list
func DeleteDish(c *gin.Context) {
	fmt.Println("Running the DeleteDish function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//GetStorageDishes gets all the dishes for the active user for a specific storage unit.
func GetStorageDishes(c *gin.Context) {
	fmt.Println("Running the GetStorageDishes function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
