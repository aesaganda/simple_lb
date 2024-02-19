package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
    "flag"
	"github.com/gorilla/mux"
)

type Server struct {
	URL        *url.URL
	ReverseURL *httputil.ReverseProxy
}

var servers []*Server
var currentIndex int
var mutex sync.Mutex

func init() {
	serverNames := []string{"localhost"}
	for i, serverName := range serverNames {
		target, _ := url.Parse(fmt.Sprintf("http://%s:%d", serverName, 8080+i))

		server := &Server{
			URL:        target,
			ReverseURL: httputil.NewSingleHostReverseProxy(target),
		}

		servers = append(servers, server)
	}
}

func getNextServer() *Server {
	mutex.Lock()
	defer mutex.Unlock()

	server := servers[currentIndex%len(servers)]
	currentIndex++

	return server
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	server := getNextServer()
	fmt.Printf("Received HTTP connection, forwarding to server at %s\n", server.URL.String())
	r.Header.Set("X-Forwarded-For", server.URL.Hostname())
	r.Header.Set("X-Forwarded-Port", server.URL.Port())
	server.ReverseURL.ServeHTTP(w, r)
}

func handleTCPConnection(clientConn net.Conn) {
	defer clientConn.Close()

	selectedServer := getNextServer()
	backendConn, err := net.Dial("tcp", selectedServer.URL.Host)
	if err != nil {
		fmt.Println("Error connecting to backend server:", err)
		return
	}
	defer backendConn.Close()

	// Print TCP Connection Information
	fmt.Printf("Received TCP connection from %s, forwarding to %s\n",
		clientConn.RemoteAddr().String(), backendConn.RemoteAddr().String())

	go func() {
		io.Copy(backendConn, clientConn)
	}()

	go func() {
		io.Copy(clientConn, backendConn)
	}()
}

func main() {
	hostname := flag.String("hostname", "0.0.0.0", "hostname to listen on")
	httpPort := flag.String("httpport", "8000", "port to listen on for HTTP")
	tcpPort := flag.String("tcpport", "9000", "port to listen on for TCP")
	flag.Parse()

	// HTTP Setup
	router := mux.NewRouter()
	router.HandleFunc("/", handleHTTPRequest).Methods(http.MethodGet)
	go func() {
		fmt.Printf("Starting HTTP load balancer on %s:%s\n", *hostname, *httpPort)
		http.ListenAndServe(*hostname+":"+*httpPort, router)
	}()

	// TCP Setup
	tcpListener, err := net.Listen("tcp", *hostname+":"+*tcpPort)
	if err != nil {
		fmt.Println("Error starting TCP listener:", err)
		return
	}
	defer tcpListener.Close()
	fmt.Printf("Starting TCP load balancer on %s:%s\n", *hostname, *tcpPort)

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			fmt.Println("Error accepting TCP connection:", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}
