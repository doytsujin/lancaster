// +build windows

package main

import (
	"net"
	"syscall"
)

func setSocketOptionInt(conn *net.UDPConn, level, option, value int) error {
	sysConn, err := conn.SyscallConn()
	if err != nil {
		return err
	}

	var serr error
	err = sysConn.Control(func(fd uintptr) {
		serr = syscall.SetsockoptInt(syscall.Handle(fd), level, option, value)
	})
	if err != nil {
		return err
	}
	return serr
}

func isENOBUFS(err error) bool {
	if err == nil {
		return false
	}

	if op, ok := err.(*net.OpError); ok {
		err = op.Err
	}
	return err == syscall.ENOBUFS
}
