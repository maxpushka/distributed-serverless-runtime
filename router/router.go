package router

import (
	"fmt"
	"log"
	"net/http"
	"serverless/router/routes_management"

	"github.com/gorilla/mux"

	"serverless/config"
	"serverless/router/auth"
	"serverless/router/database"
)

func Start(conf *config.Config) {
	db := database.Connect(conf)
	database.Initialize(db)

	router := mux.NewRouter()
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.Login(db, conf, w, r)
	}).Methods("POST")
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.Register(db, conf, w, r)
	}).Methods("POST")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(auth.Middleware(db, conf))

	api.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
		routes_management.CreateRoute(db, w, r)
	}).Methods("POST")
	api.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
		routes_management.ListRoutes(db, w, r)
	}).Methods("GET")

	api.HandleFunc("/routes/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		routes_management.GetRoute(db, w, r)
	}).Methods("GET")
	api.HandleFunc("/routes/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		routes_management.UpdateRoute(db, w, r)
	}).Methods("PUT")
	api.HandleFunc("/routes/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		routes_management.DeleteRoute(db, w, r)
	}).Methods("DELETE")

	//api.HandleFunc("/routes/{id:[0-9]+}/config", func(w http.ResponseWriter, r *http.Request) {
	//	SetConfig(db, conf, w, r)
	//}).Methods("POST")
	//api.HandleFunc("/routes/{id:[0-9]+}/config", func(w http.ResponseWriter, r *http.Request) {
	//	GetConfig(db, conf, w, r)
	//}).Methods("GET")
	//
	//api.HandleFunc("/routes/{id:[0-9]+}/executable", func(w http.ResponseWriter, r *http.Request) {
	//	SetExecutable(db, conf, w, r)
	//}).Methods("POST")
	//api.HandleFunc("/routes/{id:[0-9]+}/executable", func(w http.ResponseWriter, r *http.Request) {
	//	GetExecutable(db, conf, w, r)
	//}).Methods("GET")
	//
	//api.HandleFunc("/routes/{id:[0-9]+}/execute", func(w http.ResponseWriter, r *http.Request) {
	//	ExecuteRoute(db, conf, w, r)
	//}).Methods("POST")

	fmt.Printf("Starting server on %s\n", conf.Server.ConnectionString())
	err := http.ListenAndServe(conf.Server.ConnectionString(), router)
	if err != nil {
		log.Fatal(err)
	}

	database.Disconnect(db)
}
