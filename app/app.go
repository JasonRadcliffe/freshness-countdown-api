package app

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jasonradcliffe/freshness-countdown-api/api"
	"github.com/jasonradcliffe/freshness-countdown-api/repository/db"
	"github.com/jasonradcliffe/freshness-countdown-api/services/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/services/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/services/user"

	"github.com/gin-gonic/gin"
	domainUser "github.com/jasonradcliffe/freshness-countdown-api/domain/user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
var oauthconfig *oauth2.Config
var oauthstate string
var currentUser domainUser.OauthUser
var apiHandler api.Handler
var router = gin.Default()

func init() {
	file, err := ioutil.ReadFile("secret.config.json")
	if err != nil {
		log.Fatalln("config file error")
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalln("got an err during json.unmarshal of config" + err.Error())
	}
	oauthconfig = &oauth2.Config{
		ClientID:     config.OAuthConfig.ClientID,
		ClientSecret: config.OAuthConfig.ClientSecret,
		RedirectURL:  "https://fcapi.jasonradcliffe.com/success",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

}

//StartApplication is called by main.go and starts the app.
func StartApplication() {

	repo, err := db.NewRepository(config.DBConfig)
	if err != nil {
		log.Fatalln("StartApplication() could not create the repo")
	}

	ds := dish.NewService(repo)
	ss := storage.NewService(repo)
	us := user.NewService(repo)

	apiHandler = api.NewHandler(ds, ss, us, oauthconfig)

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

func check(err error) {
	if err != nil {
		log.Fatalln("something must have happened: ", err)
	}
}

//Privacy displays a basic privacy policy
func Privacy(c *gin.Context) {
	fmt.Println("Running the Privacy Policy function")

	siteData := []byte("<h1>Privacy Policy:</h1><br> We won't sell or send your data anywhere.<br>" +
		"Humans will review any data you submit.<br>" +
		"Your data will be kept for the purpose of maintaining and improving our service.")

	c.Data(200, "text/html", siteData)

}

func numGenerator() string {
	n := make([]byte, 32)
	rand.Read(n)
	return base64.StdEncoding.EncodeToString(n)
}
