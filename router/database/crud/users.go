package crud

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
		log.Print(err)
		return err
	}
	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", creds.Username, string(hashedPassword))
	return err
}

func GetUser(db *sql.DB, creds schema.Credentials) (schema.User, error) {
	var id int
	var username string
	var dbHashedPassword string
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", creds.Username).Scan(
		&id,
		&username,
		&dbHashedPassword,
	)
	if err != nil {
		log.Print(err)
		return schema.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbHashedPassword), []byte(creds.Password))
	if err != nil {
		return schema.User{}, errors.New("Invalid credentials")
	}
	return schema.User{
		UserId:   id,
		UserName: username,
	}, nil
}
