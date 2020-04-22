package app

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

//Config contins all the initial configuration info for this software
var config appConfig
var router = gin.Default()

func init() {
	file, err := ioutil.ReadFile("secret.config.json")
	if err != nil {
		log.Fatalln("config file error")
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("got an err during json.unmarshal of config" + err.Error())
	}


}

//StartApplication is called by main.go and starts the app.
func StartApplication() {

	mapRoutes()

	//Server Setup and Config--------------------------------------------------
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         ":443",
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		Handler:      router,
	}
	//From the config file: the file path to the fullchain .pem and privkey .pem
	log.Fatalln(srv.ListenAndServeTLS(config.CertConfig.Fullchain, config.CertConfig.PrivKey))
	//-----------------------------------------------End Server Setup and Config---

}

//Login displays a simple link that takes a user to the external google sign in flow.
func Login(c *gin.Context) {
	fmt.Println("Running the Login function")
	c.JSON(200, gin.H{
		"message": "<a href=/oauthlogin> Login with Google </a>",
	})
}


//Oauthlogin displays a simple link that takes a user to the external google sign in flow.
func Oauthlogin(c *gin.Context) {
	fmt.Println("Running the Oauthlogin function")
	c.Redirect(http.StatusTemporaryRedirect, "http://www.google.com/")
}

func Privacy(c *gin.Context){
	fmt.Println("Running the Privacy Policy function")
	c.JSON(200, gin.H{
		"message":"<h1>Privacy Policy:</h1><br> We won't sell or send your data anywhere.<br> Humans will review any data you submit.<br> Your data will be kept for the purpose of maintaining and improving our service."
	})
}