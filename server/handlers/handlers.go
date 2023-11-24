package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/koalazub/rocket-crash/templs"
)

func HandleWelcome(w http.ResponseWriter, r *http.Request) {
	err := templs.Home("chief").Render(r.Context(), w)
	if err != nil {
		slog.Error("message", "Error reading component", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleLanding(w http.ResponseWriter, r *http.Request) {
	if err := templs.Landing().Render(r.Context(), w); err != nil {
		slog.Error("couldn't make the redirect", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
}

func HandleColorchange(w http.ResponseWriter, r *http.Request) {
	newColor := "blue"

	fmt.Fprintf(w, `<button style="color: %s;" hx-trigger="click">go on</button>`, newColor)
}

func HandleRocket(w http.ResponseWriter, r *http.Request) {

}
