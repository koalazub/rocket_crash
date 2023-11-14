package server

import (
	"crypto/tls"
	"log/slog"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/quic-go/quic-go"
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
	udpConn, err := UDPServer()
	if err != nil {
		slog.Error("Couldn't intiialise UDP server", err)
	}
	InitialiseQuicTransport(udpConn)
}

func UDPServer() (*net.UDPConn, error) {
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: port})
	if err != nil {
		slog.Error("couldn't make the UDP connection at port:", err)
	}

	size := 1024 * 2048 // this is in bytes
	if err = udpConn.SetReadBuffer(size); err != nil {
		slog.Error("couldn't allocate the correct size for read buffer  ")
		return nil, err
	}
	if err = udpConn.SetWriteBuffer(size); err != nil {
		slog.Error("couldn't allocate the correct size for read buffer  ")
		return nil, err
	}

	return udpConn, nil
}

func InitialiseQuicTransport(udpConn *net.UDPConn) (*quic.Listener, error) {

	tr := quic.Transport{
		Conn: udpConn,
	}

	ln, err := tr.Listen(&tls.Config{}, &quic.Config{})
	if err != nil {
		slog.Error("Unable to listen, check tls and quic configurations", err)
		return nil, err
	}

	slog.Info("UDP Listener active on:", "", ln.Addr())

	return ln, nil
}
