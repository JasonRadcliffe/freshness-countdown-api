package storage

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//GetStorageUnits gets all the storage units for the active user.
func GetStorageUnits(c *gin.Context) {
	fmt.Println("Running the GetStorageUnits function")
	c.JSON(200, gin.H{
		"message": "Running the GetStorageUnits function",
	})
}

//GetStorageUnit gets a specific storage unit.
func GetStorageUnit(c *gin.Context) {
	storageID := c.Param("storage_id")
	fmt.Println("Running the GetStorageUnit function for the storage:", storageID)
	c.JSON(200, gin.H{
		"message": "Running the GetStorageUnit function for the storage:" + storageID,
	})
}

//CreateStorageUnit adds a storage unit to the list
func CreateStorageUnit(c *gin.Context) {
	fmt.Println("Running the CreateStorageUnit function")
	c.JSON(200, gin.H{
		"message": "Running the CreateStorageUnit function",
	})
}

//UpdateStorageUnit updates certain attributes of a specific storage unit
func UpdateStorageUnit(c *gin.Context) {
	storageID := c.Param("storage_id")
	fmt.Println("Running the UpdateStorageUnit function for the storage:", storageID)
	c.JSON(200, gin.H{
		"message": "Running the UpdateStorageUnit function for the storage:" + storageID,
	})
}

//DeleteStorageUnit deletes a specific storage unit from the list
func DeleteStorageUnit(c *gin.Context) {
	storageID := c.Param("storage_id")
	fmt.Println("Running the DeleteStorageUnit function for the storage:", storageID)
	c.JSON(200, gin.H{
		"message": "Running the DeleteStorageUnit function for the storage:" + storageID,
	})
}
