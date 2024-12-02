package routes_management

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"serverless/config"
	"serverless/router/database/crud"
	"serverless/router/schema"
)

func SetConfig(db *sql.DB, conf *config.Server, w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(schema.User)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	idStr, id, err := ParseRouteID(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(schema.Response{Error: "Invalid route ID"})
		return
	}

	baseDir, err := conf.ConfigDir()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(schema.Response{Error: "Error getting config directory"})
		return
	}

	err = SaveFile(w, r, idStr, baseDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(schema.Response{Error: "Error saving file"})
		return
	}

	err = crud.SetConfig(db, user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(schema.Response{Error: "Error setting config"})
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(schema.Response{Message: "Config set"})
}

func GetConfig(db *sql.DB, conf *config.Server, w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(schema.User)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	idStr, id, err := ParseRouteID(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(schema.Response{Error: "Invalid route ID"})
		return
	}

	_, err = crud.GetRoute(db, user, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(schema.Response{Error: "Route not found"})
		return
	}

	baseDir, err := conf.ConfigDir()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(schema.Response{Error: "Error getting config directory"})
		return
	}
	filePath, err := FindFileByName(idStr, baseDir)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	ServeFile(w, filePath)
}
