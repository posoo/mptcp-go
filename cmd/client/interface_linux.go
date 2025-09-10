//go:build linux

package main

import (
	"log"
	"net"
	"syscall"
	"unsafe"
)

// createInterfaceDialer creates a dialer that binds to a specific network interface on Linux
func createInterfaceDialer(interfaceName string) *net.Dialer {
	if interfaceName == "" {
		return &net.Dialer{}
	}

	return &net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			log.Printf("Binding to interface: %s", interfaceName)
			return c.Control(func(fd uintptr) {
				// Use SO_BINDTODEVICE to bind to specific interface
				// This is equivalent to curl's --interface option
				interfaceBytes := []byte(interfaceName)
				_, _, errno := syscall.Syscall6(
					syscall.SYS_SETSOCKOPT,
					fd,
					syscall.SOL_SOCKET,
					syscall.SO_BINDTODEVICE,
					uintptr(unsafe.Pointer(&interfaceBytes[0])),
					uintptr(len(interfaceBytes)),
					0,
				)
				if errno != 0 {
					log.Printf("Warning: Failed to bind to interface %s: %v", interfaceName, errno)
				} else {
					log.Printf("Successfully bound to interface: %s", interfaceName)
				}
			})
		},
	}
}
