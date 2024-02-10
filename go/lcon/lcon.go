package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type Server struct {
	URL               *url.URL
	ReverseURL        *httputil.ReverseProxy
	ActiveConnections int32
	Weight            int
}

var servers []*Server
var mutex sync.Mutex

func init() {
	for i := 0; i < 3; i++ {
		target, _ := url.Parse(fmt.Sprintf("http://%s:%d", "0.0.0.0", 8080+i))

		proxy := httputil.NewSingleHostReverseProxy(target)
		serverIndex := i // Create a local variable
		proxy.Director = func(req *http.Request) {
			atomic.AddInt32(&servers[serverIndex].ActiveConnections, 1) // Increment the counter here
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
		}
		proxy.ModifyResponse = func(r *http.Response) error {
			atomic.AddInt32(&servers[serverIndex].ActiveConnections, -1) // Decrement the counter here
			return nil
		}

		server := &Server{
			URL:               target,
			ReverseURL:        proxy,
			ActiveConnections: 0,
			Weight:            i + 1, // Set the weight based on the port number
		}

		servers = append(servers, server)
	}
}

func getNextServer() *Server {
	mutex.Lock()
	defer mutex.Unlock()

	// Find the server with the least number of active connections per weight
	minConnectionsPerWeight := float32(atomic.LoadInt32(&servers[0].ActiveConnections)) / float32(servers[0].Weight)
	minIndex := 0
	for i, server := range servers {
		if connectionsPerWeight := float32(atomic.LoadInt32(&server.ActiveConnections)) / float32(server.Weight); connectionsPerWeight < minConnectionsPerWeight {
			minConnectionsPerWeight = connectionsPerWeight
			minIndex = i
		}
	}

	// Increment the number of active connections for the selected server
	atomic.AddInt32(&servers[minIndex].ActiveConnections, 1)

	return servers[minIndex]
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	server := getNextServer()
	fmt.Printf("Received connection, forwarding to server at %s\n", server.URL.String())
	r.Header.Set("X-Forwarded-For", server.URL.Hostname())
	r.Header.Set("X-Forwarded-Port", server.URL.Port())

	server.ReverseURL.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8000", nil)
}
