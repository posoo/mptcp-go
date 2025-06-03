package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
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
	if err != nil {
		log.Fatal(err)
	}
	defer ls.Close()

	log.Printf("MPTCP HTTP Server listening on %s", ls.Addr())

	// Create HTTP server with custom handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHTTPRequest)

	server := &http.Server{
		Handler: mux,
	}

	// Serve HTTP requests on the MPTCP listener
	log.Fatal(server.Serve(ls))
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	remoteAddr := r.RemoteAddr
	log.Printf("[<%s] New HTTP resquest %s %s", remoteAddr, r.Method, r.URL.Path)

	hj, ok := w.(http.Hijacker)
	if !ok {
		ErrMsg := "Server doesn't support hijacking"
		http.Error(w, ErrMsg, http.StatusInternalServerError)
		log.Printf("[!%s] Error: %s", remoteAddr, ErrMsg)
		return
	}
	// Once it's hijacked, `w` is no longer usable.`
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		ErrMsg := "Failed to hijack connection: " + err.Error()
		http.Error(w, ErrMsg, http.StatusInternalServerError)
		log.Printf("[!%s] Error: %s", remoteAddr, ErrMsg)
		return
	}
	// Don't forget to close the connection when you're done.
	defer conn.Close()

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		log.Printf("[!%s] Error: Not a TCP connection. Ignore.", remoteAddr)
		return
	}

	respCode := http.StatusOK
	respErrMsg := ""

	// Access MPTCP status
	isMptcp, err := tcpConn.MultipathTCP()
	if err != nil {
		respCode = http.StatusInternalServerError
		respErrMsg = "Cannot access MPTCP status: " + err.Error()
		log.Printf("[!%s] Error: %s", remoteAddr, respErrMsg)
	}

	// Compose the response
	fmt.Fprintf(bufrw, "HTTP/1.1 %d %s\r\n", respCode, http.StatusText(respCode))
	fmt.Fprintf(bufrw, "Content-Type: text/plain\r\n")
	fmt.Fprintf(bufrw, "Date: %s\r\n", now.UTC().Format(http.TimeFormat))
	fmt.Fprintf(bufrw, "Server: MPTCP-Go-Checker\r\n")

	var respBody string
	if respCode == http.StatusOK {
		if isMptcp {
			respBody = fmt.Sprintf("You are connected via MPTCP from %s\n", remoteAddr)
		} else {
			respBody = fmt.Sprintf("You are connected via regular TCP from %s\n", remoteAddr)
		}
	} else {
		respBody = fmt.Sprintf("Error: %s\n", respErrMsg)
	}
	fmt.Fprintf(bufrw, "Content-Length: %d\r\n", len(respBody))

	// Custom headers
	fmt.Fprintf(bufrw, "X-Timestamp: %d\r\n", now.Unix())
	fmt.Fprintf(bufrw, "X-Client-Address: %s\r\n", remoteAddr)
	fmt.Fprintf(bufrw, "X-MPTCP: %t\r\n", isMptcp)
	if respErrMsg != "" {
		fmt.Fprintf(bufrw, "X-Error: %s\r\n", respErrMsg)
	}
	fmt.Fprintf(bufrw, "\r\n")
	fmt.Fprint(bufrw, respBody)

	// Flush the buffered response to the client
	bufrw.Flush()
	log.Printf("[>%s] Sent HTTP Response (MPTCP: %t)", remoteAddr, isMptcp)

}
