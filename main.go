package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

type appConfig struct {
	DBConfig   string `json:"dbCon"`
	CertConfig struct {
		Fullchain string `json:"fullchain"`
		PrivKey   string `json:"privkey"`
	} `json:"certconfigs"`
	OAuthConfig struct {
		ClientID     string `json:"clientid"`
		ClientSecret string `json:"clientsecret"`
	} `json:"oauthconfigs"`
}

func init() {
	file, err := ioutil.ReadFile("secret.config.json")
	if err != nil {
		log.Fatalln("config file error")
	}
	json.Unmarshal(file, &config)
}

var config appConfig

func main() {
	router := gin.Default()

	log.Fatalln(router.Run(":80"))

	router.GET("/ping", ping)

}

func ping(c *gin.Context) {
	fmt.Println("Running the Ping function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
