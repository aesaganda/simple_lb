package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Server struct {
	URL        *url.URL
	ReverseURL *httputil.ReverseProxy
}

var servers []*Server

func init() {
	for i := 0; i < 3; i++ {
		target, _ := url.Parse(fmt.Sprintf("http://%s:%d", "0.0.0.0", 8080+i))

		server := &Server{
			URL:        target,
			ReverseURL: httputil.NewSingleHostReverseProxy(target),
		}

		servers = append(servers, server)
	}
}

func getRandomServer() *Server {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	return servers[rand.Intn(len(servers))]
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	server := getRandomServer()
	fmt.Printf("Received connection, forwarding to server at %s\n", server.URL.String())
	r.Header.Set("X-Forwarded-For", server.URL.Hostname())
	r.Header.Set("X-Forwarded-Port", server.URL.Port())
	server.ReverseURL.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8000", nil)
}
