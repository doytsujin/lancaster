// +build darwin dragonfly freebsd linux netbsd openbsd solaris

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
		serr = syscall.SetsockoptInt(int(fd), level, option, value)
	})
	if err != nil {
		return err
	}
	return serr
}