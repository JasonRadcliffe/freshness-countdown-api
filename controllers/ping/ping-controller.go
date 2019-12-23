package ping

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//Ping is the test function to see if the server is being hit.
func Ping(c *gin.Context) {
	fmt.Println("Running the Ping function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
