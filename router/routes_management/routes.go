package routes_management

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"serverless/router/database/crud"
	"serverless/router/schema"
)

func parseRouteID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return id, nil
}

func CreateRoute(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(schema.User)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var route schema.RouteName
	err := json.NewDecoder(r.Body).Decode(&route)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(schema.Response{Error: "Invalid request body"})
		return
	}

	createdRoute, err := crud.SaveRoute(db, user, route)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(schema.Response{Error: "Error saving route"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(schema.Response{Message: "Route created", Data: *createdRoute})
}

func ListRoutes(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(schema.User)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	routes, err := crud.GetRoutes(db, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(schema.Response{Error: "Error fetching routes"})
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(schema.Response{Message: "Routes fetched", Data: routes})
}

func GetRoute(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(schema.User)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	id, err := parseRouteID(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(schema.Response{Error: "Invalid route ID"})
		return
	}

	route, err := crud.GetRoute(db, user, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(schema.Response{Error: "Route not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(schema.Response{Message: "Route fetched", Data: route})
}

func UpdateRoute(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(schema.User)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	id, err := parseRouteID(r)
	if err != nil {
		http.Error(w, "Invalid route ID", http.StatusBadRequest)
		return
	}

	var route schema.RouteName
	err = json.NewDecoder(r.Body).Decode(&route)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(schema.Response{Error: "Invalid request body"})
		return
	}

	err = crud.UpdateRoute(db, user, id, route)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(schema.Response{Error: "Error updating route"})
		return
	}

	updatedRoute, err := crud.GetRoute(db, user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(schema.Response{Error: "Error fetching updated route"})
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(schema.Response{Message: "Route updated", Data: *updatedRoute})
}

func DeleteRoute(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(schema.User)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	id, err := parseRouteID(r)
	if err != nil {
		http.Error(w, "Invalid route ID", http.StatusBadRequest)
		return
	}

	err = crud.DeleteRoute(db, user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(schema.Response{Error: "Error deleting route"})
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(schema.Response{Message: "Route deleted"})
}
