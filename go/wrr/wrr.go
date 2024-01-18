package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Server struct {
	URL        *url.URL
	ReverseURL *httputil.ReverseProxy
	Weight     int
}

var servers []*Server
var mutex sync.Mutex

func init() {
	for i := 0; i < 3; i++ {
		target, _ := url.Parse(fmt.Sprintf("http://%s:%d", "127.0.0.1", 8080+i))

		server := &Server{
			URL:        target,
			ReverseURL: httputil.NewSingleHostReverseProxy(target),
			Weight:     i + 1, // Assigning weights. Modify as per your requirements.
		}

		servers = append(servers, server)
	}
}

func getNextServer() *Server {
	mutex.Lock()
	defer mutex.Unlock()

	// Calculate the total weight of all servers
	total := 0
	for _, server := range servers {
		total += server.Weight
	}

	// Generate a random number between 0 and the total weight
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(total)

	// Select the server that corresponds to the random number
	for _, server := range servers {
		if random < server.Weight {
			return server
		}
		random -= server.Weight
	}

	// Fallback to the first server if no server is selected (this should never happen)
	return servers[0]
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
	fmt.Println("Starting server on port 8000")
	http.ListenAndServe(":8000", nil)
}
