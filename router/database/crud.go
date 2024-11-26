package database

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"serverless/router/schema"
)

func SaveUser(db *sql.DB, creds schema.Credentials) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", creds.Username, string(hashedPassword))
	return err
}

func GetUserPassword(db *sql.DB, creds schema.Credentials) error {
	var dbHashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", creds.Username).Scan(&dbHashedPassword)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbHashedPassword), []byte(creds.Password))
	if err != nil {
		return errors.New("Invalid credentials")
	}
	return nil
}
