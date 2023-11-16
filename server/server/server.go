package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/quic-go/quic-go/http3"
)

func init() {
	loadEnv()
}

var addr string
var port int

func loadEnv() {
	// Reload env later on via cicd - this looks like shit
	if err := godotenv.Load("../.env"); err != nil {
		slog.Error("Couldn't load env variables for server. Are they there? ")
		return
	}

	addr = os.Getenv("server_addr")
	t := os.Getenv("server_port")
	p, err := strconv.Atoi(t)
	if err != nil {
		slog.Error("Couldn't get port:", err)
		return
	}
	port = p
}

func InitServer() {
	fAddr := addr + ":" + string(port)
	crt := "certfile.cert"
	key := "keyfile.key"
	slog.Info("Now listening on:", fAddr)
	err := http3.ListenAndServeQUIC(fAddr, crt, key, setupHandler())
	if err != nil {
		slog.Error("Couldn't create QUIC server. Check that the address, cert, key and handlers are good: ", err)
	}

}

func setupHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/rocket", handleRocket)

	return mux
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "we is in da home")
}
func handleRocket(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "this is where the rockets should be rendered")
}
