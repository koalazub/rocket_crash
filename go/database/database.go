package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var tursoAuth, tursoUrl, dbUrl string

func init() {
	loadEnv()
}

func StartDatabase() *sql.DB {
	db := initConnection()
	err := initTable(db)
	if err != nil {
		slog.Error("Couldn't start database", err)
		return nil
	}

	slog.Info("Started database", "", db.Stats())
	return db
}

func initConnection() *sql.DB {
	databaseUrl := fmt.Sprintf("libsql://%s?authToken=%s", tursoUrl, tursoAuth)
	db, err := sql.Open("libsql", databaseUrl)
	if err != nil {
		slog.Error("database couldn't be opened", "Err:", err)
		os.Exit(1)
	}

	return db
}

func initTable(db *sql.DB) error {
	_, err := db.Exec("create table if not exists rockets(id INT, name varchar(255),  crashed int, death_coord_x real, death_coord_y real, rocket_type varchar(255))")
	if err != nil {
		slog.Error("Error creating table. Verify that the query is correct")
		return err
	}
	slog.Info("rocket table initialised")
	return nil
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		{
			slog.Error("Error loading db vars. Are they present? ", err)
			return
		}
	}

	tursoUrl = os.Getenv("turso_url")
	tursoAuth = os.Getenv("turso_auth")
	if tursoAuth == "" || tursoUrl == "" {
		slog.Error("couldn't load env variables. Are they there?")
		return
	}

}
