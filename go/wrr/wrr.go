package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	URL        *url.URL
	ReverseURL *httputil.ReverseProxy
	Weight     int
}

var servers []*Server
var currentIndex = 0
var currentWeight = 0

func init() {
	for i := 0; i < 3; i++ {
		target, _ := url.Parse(fmt.Sprintf("http://%s:%d", "0.0.0.0", 8080+i))

		server := &Server{
			URL:        target,
			ReverseURL: httputil.NewSingleHostReverseProxy(target),
			Weight:     i + 1, // Set the weight based on the port number
		}

		servers = append(servers, server)
	}
}

func getNextServer() *Server {
	server := servers[currentIndex]
	if currentWeight >= server.Weight {
		currentWeight = 0
		currentIndex = (currentIndex + 1) % len(servers)
	}
	currentWeight++
	return server
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
