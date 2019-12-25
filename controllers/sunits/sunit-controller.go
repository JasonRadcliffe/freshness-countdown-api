package sunits

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//GetSUnits gets all the storage units for the active user.
func GetSUnits(c *gin.Context) {
	fmt.Println("Running the GetSUnits function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//GetSUnit gets a specific storage unit.
func GetSUnit(c *gin.Context) {
	fmt.Println("Running the GetSUnit function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//CreateSUnit adds a storage unit to the list
func CreateSUnit(c *gin.Context) {
	fmt.Println("Running the CreateSUnit function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//UpdateSUnit updates certain attributes of a specific storage unit
func UpdateSUnit(c *gin.Context) {
	fmt.Println("Running the UpdateSUnit function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//DeleteSUnit deletes a specific storage unit from the list
func DeleteSUnit(c *gin.Context) {
	fmt.Println("Running the DeleteSUnit function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
