package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
)

func main() {
	hostname := flag.String("hostname", "0.0.0.0", "hostname to listen on")
	httpPort := flag.String("httpport", "8080", "port to listen on for HTTP")
	tcpPort := flag.String("tcpport", "9090", "port to listen on for TCP")
	flag.Parse()

	// Start HTTP server in a goroutine
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			forwardedHostname := r.Header.Get("X-Forwarded-For")
			forwardedPort := r.Header.Get("X-Forwarded-Port")

			// Print received HTTP connection information
			fmt.Printf("Received HTTP connection from %s:%s\n", forwardedHostname, forwardedPort)

			fmt.Fprintf(w, "Hello from server at Hostname: %s and Port: %s\n", forwardedHostname, forwardedPort)
		})
		fmt.Printf("Starting HTTP server on %s:%s\n", *hostname, *httpPort)
		http.ListenAndServe(*hostname+":"+*httpPort, nil)
	}()

	// Start TCP listener in a goroutine
	go func() {
		listener, err := net.Listen("tcp", *hostname+":"+*tcpPort)
		if err != nil {
			fmt.Println("Error starting TCP listener:", err)
			return
		}
		defer listener.Close()
		fmt.Printf("Starting TCP server on %s:%s\n", *hostname, *tcpPort)

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting TCP connection:", err)
				continue
			}
			go handleTCPConnection(conn)
		}
	}()

	// Keep the main goroutine alive (optional)
	select {}
}

func handleTCPConnection(conn net.Conn) {
    defer conn.Close()
    fmt.Printf("Received TCP connection from %s\n", conn.RemoteAddr().String())

    // Data handling loop
    for {
        // Read data from the client connection
        buffer := make([]byte, 1024)  // Example buffer size
        n, err := conn.Read(buffer)
        if err != nil {
            if err != io.EOF { //  Handle non-EOF errors
                fmt.Println("Error reading from connection:", err)
            }
            break // Connection closed (either error or intentionally by the client)
        }

        // Process the received data (do something with buffer[:n]) 
        fmt.Println("Received data:", string(buffer[:n])) 

        // Optionally send a response back to the client 
        // conn.Write(...) 
    }
}
