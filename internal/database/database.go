package DB

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var Db *sql.DB

func SetDbConfig() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.Ping()
}

func InsertUser(username, weight, height, gender string, lost, set, get int) error {
	_, err := Db.Exec(`
		INSERT INTO bot_users (username, weight, height, gender, Kforlost, Kforset, Kforget) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		username, weight, height, gender, lost, get, set)
	return err
}
