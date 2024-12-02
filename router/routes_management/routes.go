package routes_management

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"serverless/router/database/crud"
	"serverless/router/schema"
)

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
