package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var tursoAuth, tursoUrl, dbUrl string
var ToLog *bool

func init() {
	loadEnv()
}

// don't json this - using capnp
type Rocket struct {
	ID          int64
	Name        string
	User        string
	Crashed     bool
	RocketType  string
	deathCoordX float64
	deathCoordY float64
}

func Start() *sql.DB {
	db := initConnection()
	err := initTable(db)
	if err != nil {
		slog.Error("Couldn't start database", err)
		return nil
	}

	if ToLog != nil && *ToLog {
		slog.Info("Started database", "", db.Stats())
	}
	return db
}

func GetRockets(db *sql.DB) ([]Rocket, error) {
	rows, err := db.Query("select id, name, crashed, death_coord_x, death_coord_y, rocket_type from rockets")
	if err != nil {
		slog.Error("couldn't fetch from the database")
		return nil, err
	}

	defer rows.Close()
	var rockets []Rocket
	for rows.Next() {
		var r Rocket
		err := rows.Scan(&r.ID, &r.Name, &r.User, &r.Crashed, &r.RocketType, &r.deathCoordX, &r.deathCoordY)
		if err != nil {
			slog.Error("Error scanning rocket rows ", "Err:", err)
			return nil, err
		}
		rockets = append(rockets, r)
	}
	if ToLog != nil && *ToLog {
		slog.Info("query", "rows: ", rows)
		slog.Info("query", "rockets: ", rockets)
	}
	return rockets, nil

}
func initConnection() *sql.DB {
	databaseUrl := fmt.Sprintf("libsql://%s?authToken=%s", tursoUrl, tursoAuth)
	db, err := sql.Open("libsql", databaseUrl)
	if err != nil {
		slog.Error("database couldn't be opened", "Err:", err)
		os.Exit(1)
	}

	if ToLog != nil && *ToLog {
		slog.Info("database", "url", databaseUrl)
		slog.Info("database init connection", "db", db)
	}

	return db
}

var createTableSQL = `
		CREATE TABLE IF NOT EXISTS rockets (
		id INT, 
		name varchar(255)
		crashed INT,
		death_coord_x real,
		death_coord_y real,
		rocket_type varchar(255) 
	)
`

func initTable(db *sql.DB) error {
	requiredCols := map[string]string{
		"id":            "INT",
		"name":          "varchar(255)",
		"crashed":       "INT",
		"death_coord_x": "real",
		"death_coord_y": "real",
		"rocket_type":   "varchar(255)",
	}

	_, err := db.Exec(createTableSQL)
	if err != nil {
		slog.Error("Error creating table. Verify that the query is correct")
		return err
	}
	slog.Info("rocket table initialised")

	for colName, colType := range requiredCols {
		var cn string
		err := db.QueryRow(`
				SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS 
				WHERE TABLE_NAME = 'rockets' AND COLUMN_NAME = ? 
			`, cn).Scan(&colName)
		if err == sql.ErrNoRows {
			alterSQL := fmt.Sprintf("ALTER TABLE rockets ADD COLUMN %s %s", colName, colType)
			_, err = db.Exec(alterSQL)
			if err != nil {
				slog.Error("Error adding column", "Column: ", err)
				return err
			}

			slog.Info("Column added", "column: ", colName)
		} else if err != nil {
			slog.Error("Error checking for column", "column: ", err)
			return err
		}

	}
	slog.Info("rocket table's been initialised")
	return nil
}

func loadEnv() {
	if err := godotenv.Load("../.env"); err != nil {
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
