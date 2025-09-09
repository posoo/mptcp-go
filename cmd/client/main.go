package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	noMptcp := flag.Bool("no-mptcp", false, "Use regular TCP instead of MPTCP")
	flag.Parse()
	args := flag.Args()
	// Get remote address and message count from command line
	if len(args) != 2 {
		fmt.Println("Usage: client [--no-mptcp] <remote_address> <message_count>")
		os.Exit(1)
	}
	fmt.Println(args)
	remoteAddr := args[0]
	messageCount, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Failed to parse message count:", err)
		os.Exit(1)
	}

	// Initiate MPTCP connection
	dialer := &net.Dialer{}
	if *noMptcp {
		log.Println("Creating single TCP connection...")
	} else {
		log.Println("Creating single MPTCP connection...")
		dialer.SetMultipathTCP(true)
	}
	conn, err := dialer.Dial("tcp", remoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Log received data
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Received %d bytes: %s", n, string(buffer[:n]))
		}
	}()

	// Send data over the MPTCP connection
	// Traffic should automatically use both subflows if they exist
	for i := 0; i < messageCount; i++ {
		message := fmt.Sprintf("Hello, world! Message %d\n", i)
		conn.Write([]byte(message))
		// time.Sleep(100 * time.Millisecond)
		if i%10 == 0 {
			time.Sleep(1 * time.Second)
		}
		log.Printf("Sent message %d\n", i)
	}

	// Wait for all messages to be received
	time.Sleep(2 * time.Second)

	log.Printf("Done")
}
