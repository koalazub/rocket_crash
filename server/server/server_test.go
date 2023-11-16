package server

import (
	"log/slog"
	"net"
	"syscall"
)

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
