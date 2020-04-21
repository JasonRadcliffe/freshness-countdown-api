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
	err2 := json.Unmarshal(file, &config)
	if err2 != nil {
		fmt.Println("got an err during json.unmarshal" + err2.Error())
	}
	fmt.Println(config)
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
	fmt.Println("about to listen and serve tls" + config.CertConfig.Fullchain)
	log.Fatalln(srv.ListenAndServeTLS(config.CertConfig.Fullchain, config.CertConfig.PrivKey))
	//-----------------------------------------------End Server Setup and Config---

}
