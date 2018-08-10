package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/getyoti/yoti-go-sdk"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// TODO: for deploy, those arguments should get from os.Getenv()
	appID  = "<Yoti AppID>"
	sdkID  = "<Yoti SDKID>"
	dbConn = "<MySQL Connection>"
	port   = "8080"

	sdkKeyFile    = "key.pem"
	htmlTemplte   = "index.html"
	anonymousFile = "anonymous.png"

	sessionID         = "yoti"
	sessionProfileKey = "profile"
)

var store = newCookieStore()

func yotiClient(skdID, keyFile string) (*yoti.Client, error) {

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	client := &yoti.Client{
		SdkID: skdID,
		Key:   key,
	}

	return client, nil
}

func main() {

	// TODO: Log error to os.StdErr, no panic
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	c, err := yotiClient(sdkID, sdkKeyFile)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Handle("/", mainHandler(db))
	r.Handle("/comments", postHandler(db))
	r.Handle("/callback", callbackHandler(c, db))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler(db, anonymousFile)))

	fmt.Fprintf(os.Stdout, "Try listening...:%s", port)
	if err := http.ListenAndServe(":"+port, recoveryHandler(true, r)); err != nil {
		panic(err)
	}
}
