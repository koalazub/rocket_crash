package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	d "github.com/koalazub/rocket-crash/database"
	"github.com/koalazub/rocket-crash/templs"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Rando"
	}
	err := templs.Home(name).Render(r.Context(), w)
	if err != nil {
		slog.Error("message", "Error reading component", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Landing(w http.ResponseWriter, r *http.Request) {
	if err := templs.Landing().Render(r.Context(), w); err != nil {
		slog.Error("couldn't make the redirect", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
}

func Rockets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rockets, err := d.GetRockets(db)
		if err != nil {
			slog.Error("Couldn't fetch rockets", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		for _, r := range rockets {
			fmt.Printf("Rockets! %v", r)
		}
	}
}
