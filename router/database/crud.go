package database

import (
	"database/sql"
	"serverless/router/schema"
)

func SaveUser(db *sql.DB, creds schema.Credentials) error {
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", creds.Username, creds.Password)
	return err
}

func GetUserPassword(db *sql.DB, creds schema.Credentials) (string, error) {
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", creds.Username).Scan(&dbPassword)
	return dbPassword, err
}
