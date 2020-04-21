package storage

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//GetStorageUnits gets all the storage units for the active user.
func GetStorageUnits(c *gin.Context) {
	fmt.Println("Running the GetStorageUnits function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//GetStorageUnit gets a specific storage unit.
func GetStorageUnit(c *gin.Context) {
	fmt.Println("Running the GetStorageUnit function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//CreateStorageUnit adds a storage unit to the list
func CreateStorageUnit(c *gin.Context) {
	fmt.Println("Running the CreateStorageUnit function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//UpdateStorageUnit updates certain attributes of a specific storage unit
func UpdateStorageUnit(c *gin.Context) {
	fmt.Println("Running the UpdateStorageUnit function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//DeleteStorageUnit deletes a specific storage unit from the list
func DeleteStorageUnit(c *gin.Context) {
	fmt.Println("Running the DeleteStorageUnit function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
