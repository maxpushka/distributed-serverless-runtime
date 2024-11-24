package auth

import (
	"database/sql"
	"github.com/golang-jwt/jwt"
	"net/http"
	"serverless/config"
	"time"
)

func Login(db *sql.DB, conf *config.Config, w http.ResponseWriter, r *http.Request) {
	// Parse the request body for username and password
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate inputs
	if username == "" || password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	// Check user credentials
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&dbPassword)
	if err != nil || dbPassword != password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(conf.AuthJWTExpires)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(conf.AuthJWTKey)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	// Return token
	w.Write([]byte(tokenStr))
}

func Register(db *sql.DB, conf *config.Config, w http.ResponseWriter, r *http.Request) {
	// Parse the request body for username and password
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate inputs
	if username == "" || password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	// Save user to the database
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User registered successfully"))
}
