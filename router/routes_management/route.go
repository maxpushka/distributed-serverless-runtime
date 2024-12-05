package routes_management

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"serverless/executor"
	"serverless/executor/js"

	"serverless/router/database/crud"
	"serverless/router/schema"
)

func GetRoute(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(schema.User)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	_, id, err := ParseRouteID(r)
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

	_, id, err := ParseRouteID(r)
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

	_, id, err := ParseRouteID(r)
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

func ExecuteRoute(db *sql.DB, executorJs *js.Executor, w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(schema.User)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	idStr, id, err := ParseRouteID(r)
	if err != nil {
		http.Error(w, "Invalid route ID", http.StatusBadRequest)
		return
	}

	var params schema.ExecuteParams
	err = json.NewDecoder(r.Body).Decode(&params)

	route, err := crud.GetRoute(db, user, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(schema.Response{Error: "Route not found"})
		return
	}

	if !(route.ExecutableExists) {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(schema.Response{Error: "Route not executable"})
		return
	}

	request := executor.Request{
		Method:  "POST",
		URL:     idStr,
		Headers: nil,
		Body:    params.RequestBody,
	}

	response, err := executorJs.Execute(idStr, params.RouteConfig, request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(schema.Response{Error: "Error executing route"})
		return
	}
	w.WriteHeader(http.StatusOK)
	encoder.Encode(schema.Response{Message: "Route executed", Data: response})
}
