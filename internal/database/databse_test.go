package database

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func Test_dbconnect(T *testing.T) {

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	_, err := sql.Open("postgres", connStr)
	if err != nil {
		T.Errorf("postgres is not connected")
	}
}
