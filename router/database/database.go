package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"serverless/config"
)

func Connect(conf *config.Config) (*sql.DB, error) {
	db, errOpen := sql.Open("postgres", conf.DbUrl())
	if errOpen != nil {
		log.Fatal(errOpen)
		return db, errOpen
	}
	_, errTestConnection := db.Query("SELECT 1")
	if errTestConnection != nil {
		log.Fatal(errTestConnection)
		return db, errTestConnection
	}
	return db, nil
}

func Initialize(db *sql.DB) error {
	// TODO(Vlad): Create users and configs table here
	return nil
}

func Close(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
	return err
}
