package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	h "github.com/koalazub/rocket-crash/handlers"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/logging"
	"github.com/quic-go/quic-go/qlog"
)

func init() {
	loadEnv()
}

var addr string
var port int

func loadEnv() {
	// Reload env later on via cicd - this looks like shit
	if err := godotenv.Load("../.env"); err != nil {
		slog.Error("Couldn't load env variables for server. Are they there?", "\n", err)
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
	initLogger()
	fAddr := addr + ":" + fmt.Sprintf("%d", port)
	crt := "certfile.crt"
	key := "keyfile.key"

	slog.Info("message", "Now listening on: ", fAddr)

	err := http3.ListenAndServeQUIC(fAddr, crt, key, setupHandler())

	if err != nil {
		slog.Error("ListenAndServeQuic error:  ", "error", err.Error())
		return
	}

}

func initLogger() *quic.Config {
	slog.Info("logger initialised")
	return &quic.Config{
		Tracer: func(ctx context.Context, p logging.Perspective, ci quic.ConnectionID) *logging.ConnectionTracer {
			role := "server"
			if p == logging.PerspectiveClient {
				role = "client"
			}
			filename := fmt.Sprintf("./log_%s_%s.qlog", ci, role)
			f, err := os.Create(filename)
			if err != nil {
				slog.Error("Error during quic log file creation process: ", "error ", err.Error())
				return nil
			}

			return qlog.NewConnectionTracer(f, p, ci)
		},
	}
}

func setupHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HandleLanding)
	mux.HandleFunc("/welcome", h.HandleWelcome)
	mux.HandleFunc("/rocket", h.HandleRocket)

	return mux
}
