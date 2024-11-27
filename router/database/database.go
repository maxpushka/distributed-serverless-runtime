package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"serverless/config"
)

func Connect(conf *config.Config) *sql.DB {
	db, errOpen := sql.Open("postgres", conf.Db.ConnectionString())
	if errOpen != nil {
		log.Fatal(errOpen)
	}
	_, errTestConnection := db.Query("SELECT 1")
	if errTestConnection != nil {
		log.Fatal(errTestConnection)
	}
	return db
}

func Initialize(db *sql.DB) {
	initialization := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE,
			password TEXT
		);
		CREATE TABLE IF NOT EXISTS routes (
			id SERIAL PRIMARY KEY,
			name TEXT,
			config_exists BOOLEAN DEFAULT FALSE,
			executable_exists BOOLEAN DEFAULT FALSE,
			user_id INTEGER REFERENCES users(id)
		);
	`
	_, err := db.Exec(initialization)
	if err != nil {
		log.Fatal(err)
	}
}

func Disconnect(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
