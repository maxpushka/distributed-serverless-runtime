package crud

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"

	"serverless/router/schema"
)

func SaveUser(db *sql.DB, creds schema.Credentials) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return err
	}
	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", creds.Username, string(hashedPassword))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func CheckUserPassword(db *sql.DB, creds schema.Credentials) bool {
	var dbHashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", creds.Username).Scan(&dbHashedPassword)
	if err != nil {
		log.Print(err)
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbHashedPassword), []byte(creds.Password))
	if err != nil {
		log.Print(err)
		return false
	}
	return true
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
	return schema.User{
		UserId:   id,
		UserName: username,
	}, nil
}
