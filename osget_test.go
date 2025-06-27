package main

import (
	"os"
	"testing"
)

func Test_OsGet(T *testing.T) {
	slice := []string{os.Getenv("DB_USER"), os.Getenv("DB_PORT"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("TELEGRAM_BOT_TOKEN")}
	for i, j := range slice {
		if j == "" {
			T.Error("Not correct os data", i)
		}
	}
}
