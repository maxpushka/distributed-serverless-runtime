package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"serverless/router/database/crud"
	"time"

	"github.com/golang-jwt/jwt"

	"serverless/config"
	"serverless/router/schema"
)

func Login(db *sql.DB, conf *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	// Parse the request body for username and password
	var creds schema.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(schema.Response{Error: "Invalid request body"})
		return
	}

	// Validate inputs
	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(schema.Response{Error: "Username and password required"})
		return
	}

	// Check user credentials
	user, err := crud.GetUser(db, creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(schema.Response{Error: "Invalid credentials"})
		return
	}

	// Create JWT token
	issuedAt := time.Now()
	expirationTime := issuedAt.Add(conf.Auth.JWTExpires)
	claims := user.ToClaims(jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  issuedAt.Unix(),
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(conf.Auth.JWTKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(schema.Response{Error: "Failed to create token"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schema.Response{Message: "Login successful", Data: schema.TokenData{Token: tokenStr}})
}

func Register(db *sql.DB, conf *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var creds schema.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(schema.Response{Error: "Invalid request body"})
		return
	}

	// Validate inputs
	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(schema.Response{Error: "Username and password required"})
		return
	}

	// Save user to the database
	err = crud.SaveUser(db, creds)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		encoder.Encode(schema.Response{Error: "User with this username already exists"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schema.Response{Message: "User created successfully"})
}
