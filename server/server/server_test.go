package server

import (
	"log/slog"
	"net"
	"syscall"
	"testing"
)

func TestInitialiseUDPServer(t *testing.T) {
	got := UDPServer()
	var want *net.UDPConn = &net.UDPConn{}

	if got != want {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}

func TestUDPBufferSize(t *testing.T) {
	udpSrv := UDPServer()
	var srv *net.UDPConn = &net.UDPConn{}
	srv.SetReadBuffer(1024 * 2048)
	srv.SetWriteBuffer(1024 * 2048)

	gotBufSize := receiveBufSize(udpSrv)
	wantBufSize := receiveBufSize(udpSrv)

	if gotBufSize != wantBufSize {
		t.Errorf("got: %v, want: %v", gotBufSize, wantBufSize)
	}
}

func receiveBufSize(s *net.UDPConn) int {
	var rcvBufSize int
	sysCall, err := s.SyscallConn()
	if err != nil {
		slog.Error("couldn't make raw connection to system call", err)
	}
	sysCall.Control(func(fd uintptr) {
		rcvBufSize, err = syscall.GetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF)
		if err != nil {
			slog.Error("Unable to read buffer size from system call", err)
		}
	})

	return rcvBufSize
}
