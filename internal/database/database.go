package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func SetDbConfig() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("connection error")
	}

	log.Println("postgres connect")
	return nil
}

func InsertUser(username, weight, height, gender, program string, lost, set, get int) error {
	if Db == nil {
		return fmt.Errorf("error db")
	}

	_, err := Db.Exec(`
        INSERT INTO bot_users (username, weight, height, gender, kforlost, kforset, kforget, program) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (username) 
        DO UPDATE SET 
            weight = EXCLUDED.weight,
            height = EXCLUDED.height,
            gender = EXCLUDED.gender,
            kforlost = EXCLUDED.kforlost,
            kforset = EXCLUDED.kforset,
            kforget = EXCLUDED.kforget,
			program = EXCLUDED.program`,

		username, weight, height, gender, lost, set, get, program)

	return err
}
