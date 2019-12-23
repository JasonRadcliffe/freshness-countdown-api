package main

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

	router.GET("/ping", ping)

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

func ping(c *gin.Context) {
	fmt.Println("Running the Ping function: PONG")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
