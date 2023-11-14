package server

import (
	"log/slog"
	"net"
	"syscall"
	"testing"
)

func TestInitialiseUDPServer(t *testing.T) {
	conn, err := UDPServer()
	if err != nil {
		t.Fatalf("UDPServer() failed with error: %v", err)
	}
	if conn == nil {
		t.Fatalf("UDPServer() return nil connection")
	}
}

// Buffer size varies according to the OS by up to twice as much
func TestUDPBufferSize(t *testing.T) {
	udpSrv, err := UDPServer()
	if err != nil {
		t.Fatalf("UDPServer9) failed with error: %v", err)
	}

	got := receiveBufSize(udpSrv)
	minwant := 1024 * 2048 // 2048KiB
	maxwant := 1024 * 4096 // 4096KiB
	if got < minwant || got > maxwant {
		t.Fatalf("got: %v, want between: %v - %v", got, minwant, maxwant)
	}
}

func TestQuicTransport(t *testing.T) {
	conn, _ := UDPServer()
	qt, err := InitialiseQuicTransport(conn)
	if err != nil {
		t.Fatalf("Error invoking InitialiseQuicTransport(): %v", err)
	}

	if qt == nil {
		t.Fatalf("InitialiseQuicTransport() is nil ")
	}

}
func receiveBufSize(s *net.UDPConn) int {
	var rcvBufSize int
	sysCall, err := s.SyscallConn()
	if err != nil {
		slog.Error("couldn't make raw connection to system call", err)
		return 0
	}
	sysCall.Control(func(fd uintptr) {
		rcvBufSize, err = syscall.GetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF)
		if err != nil {
			slog.Error("unable to read buffer size from system call", err)
		}
	})
	if err != nil {
		slog.Error("Error in syscall.Control", err)
		return 0
	}

	return rcvBufSize
}
