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
		create table if not exists rockets (
		id int, 
		name text,
		crashed int,
		death_coord_x real,
		death_coord_y real,
		rocket_type text 
	)
`

func initTable(db *sql.DB) error {
	requiredCols := map[string]string{
		"id":            "INT",
		"name":          "text",
		"crashed":       "INT",
		"death_coord_x": "real",
		"death_coord_y": "real",
		"rocket_type":   "text",
	}

	_, err := db.Exec(createTableSQL)
	if err != nil {
		slog.Error("Error creating table. Verify that the query is correct")
		return err
	}
	slog.Info("rocket table initialised")

	rows, err := db.Query("PRAGMA table_info(rockets);")
	if err != nil {
		slog.Error("couldn't query rockets", err)
		return err
	}

	defer rows.Close()

	existingCols := make(map[string]bool)
	for rows.Next() {
		var (
			cid      string
			name     string
			colType  string
			notnull  int
			dflt_val *string
			pk       int
		)
		if err := rows.Scan(&cid, &name, &colType, &notnull, &dflt_val, &pk); err != nil {
			slog.Error("Error scanning table: ", err)
			return err
		}
		existingCols[name] = true
	}

	for colName, colType := range requiredCols {
		if !existingCols[colName] {
			alterSQL := fmt.Sprintf("ALTER TABLE rockets ADD COLUMN %s %s", colName, colType)
			if _, err = db.Exec(alterSQL); err != nil {
				slog.Error("Couldn't alter table", err)
				return err
			}

			if ToLog != nil && *ToLog {
				slog.Info("Column added", "col:", colName)
			}
		}
	}

	if ToLog != nil && *ToLog {
		slog.Info("query", "after:", err)
		slog.Info("rocket table's been initialised")
	}

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
