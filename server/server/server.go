package server

import (
	"context"
	"database/sql"
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

var ToLog *bool

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
	if ToLog != nil && *ToLog {
		slog.Info("port information:", "data:", port)
	}
}

func InitServer(toLog *bool, db *sql.DB) {
	defer initLogger()
	// initLogger()
	fAddr := addr + ":" + fmt.Sprintf("%d", port)
	crt := "certfile.crt"
	key := "keyfile.key"

	slog.Info("message", "Now listening on: ", fAddr)

	err := http3.ListenAndServeQUIC(fAddr, crt, key, setupHandler(db))

	if err != nil {
		slog.Error("ListenAndServeQuic error:  ", "error", err.Error())
		return
	}

	if ToLog != nil && *ToLog {
		slog.Info("address info:", "faddr is: ", fAddr)
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

			if ToLog != nil && *ToLog {
				qlog.NewConnectionTracer(f, p, ci)
				slog.Info("filename saved is", "file:", filename)
			}

			return qlog.NewConnectionTracer(f, p, ci)
		},
	}
}

func setupHandler(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Landing)
	mux.HandleFunc("/welcome", h.Welcome)
	mux.HandleFunc("/rockets", h.Rockets(db))
	if ToLog != nil && *ToLog {
		fmt.Printf("mux: %v\n", mux)
	}

	return mux
}
