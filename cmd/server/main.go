package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	// Get port from command line
	if len(os.Args) != 2 {
		fmt.Println("Usage: server <port>")
		os.Exit(1)
	}
	port := os.Args[1]
	if _, err := strconv.Atoi(port); err != nil {
		fmt.Println("Usage: server <port>")
		os.Exit(1)
	}

	// Set up a MPTCP listener
	lc := &net.ListenConfig{}
	lc.SetMultipathTCP(true)
	ls, err := lc.Listen(context.Background(), "tcp", ":"+port)
	log.Println("MPTCP Server listening on", ls.Addr())
	if err != nil {
		log.Fatal(err)
	}
	defer ls.Close()

	// Accept connections
	for {
		conn, err := ls.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

// Echo server
func handleConnection(conn net.Conn) {
	log.Printf("[%s] New connection from %s\n", time.Now().Format("15:04:05"), conn.RemoteAddr())

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		log.Printf("[%s] Received %d bytes\n", time.Now().Format("15:04:05"), n)
		if err != nil {
			if err == io.EOF {
				log.Printf("[%s] Connection closed by %s\n", time.Now().Format("15:04:05"), conn.RemoteAddr())
			}
			break
		}

		conn.Write(buffer[:n])
		log.Printf("[%s] Echoed %d bytes to %s\n", time.Now().Format("15:04:05"), n, conn.RemoteAddr())
	}

	conn.Close()
}
