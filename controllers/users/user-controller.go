package users

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//GetUsers gets all the dishes for the active user.
func GetUsers(c *gin.Context) {
	fmt.Println("Running the GetUsers function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//GetUser gets a specific dish.
func GetUser(c *gin.Context) {
	fmt.Println("Running the GetUser function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//CreateUser adds a user to the list
func CreateUser(c *gin.Context) {
	fmt.Println("Running the CreateUser function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//UpdateUser updates certain attributes of a specific user
func UpdateUser(c *gin.Context) {
	fmt.Println("Running the UpdateUser function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//DeleteUser deletes a specific user from the list
func DeleteUser(c *gin.Context) {
	fmt.Println("Running the DeleteUser function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
