package router

import (
	//"context"
	//"database/sql"
	//"encoding/json"
	"fmt"
	//"io"
	"log"
	"net/http"
	//"os"
	"serverless/router/database"
	//"strings"
	//"time"

	"github.com/gorilla/mux"

	"serverless/config"
	"serverless/router/auth"
)

func Start(conf *config.Config) error {
	db, err := database.Connect(conf)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = database.Initialize(db)
	if err != nil {
		log.Fatal(err)
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.Login(db, conf, w, r)
	}).Methods("POST")
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.Register(db, conf, w, r)
	}).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.Middleware(conf))

	fmt.Printf("Starting server on port %s\n", conf.ServerPort)
	err = http.ListenAndServe(":"+conf.ServerPort, r)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = database.Close(db)
	return err
}
