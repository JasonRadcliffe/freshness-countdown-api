package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//GetUsers gets all the dishes for the active user.
func GetUsers(c *gin.Context) {
	fmt.Println("Running the GetUsers function")
	c.JSON(200, gin.H{
		"message": "Running the GetUsers function",
	})
}

//GetUser gets a specific dish.
func GetUser(c *gin.Context) {
	userID := c.Param("user_id")
	fmt.Println("Running the GetUser function for this user:", userID)
	c.JSON(200, gin.H{
		"message": "Running the GetUser function for this user:" + userID,
	})
}

//CreateUser adds a user to the list
func CreateUser(c *gin.Context) {
	fmt.Println("Running the CreateUser function")
	c.JSON(200, gin.H{
		"message": "Running the CreateUser function",
	})
}

//UpdateUser updates certain attributes of a specific user
func UpdateUser(c *gin.Context) {
	fmt.Println("Running the UpdateUser function")
	c.JSON(200, gin.H{
		"message": "Running the UpdateUser function",
	})
}

//DeleteUser deletes a specific user from the list
func DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	fmt.Println("Running the DeleteUser function for this user:", userID)
	c.JSON(200, gin.H{
		"message": "Running the DeleteUser function for this user:" + userID,
	})
}
