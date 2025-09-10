//go:build !linux

package main

import (
	"log"
	"net"
)

// createInterfaceDialer creates a dialer that ignores interface binding on non-Linux platforms
func createInterfaceDialer(interfaceName string) *net.Dialer {
	if interfaceName != "" {
		log.Printf("Warning: Interface binding (--interface) is not supported on this platform. Ignoring interface: %s", interfaceName)
	}
	return &net.Dialer{}
}
