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
    URL                *url.URL
    ReverseURL         *httputil.ReverseProxy
    ActiveConnections  int32
}

var servers []*Server
var mutex sync.Mutex

func init() {
    for i := 0; i < 3; i++ {
        target, _ := url.Parse(fmt.Sprintf("http://%s:%d", "127.0.0.1", 8080+i))

        server := &Server{
            URL:               target,
            ReverseURL:        httputil.NewSingleHostReverseProxy(target),
            ActiveConnections: 0,
        }

        servers = append(servers, server)
    }
}

func getNextServer() *Server {
    mutex.Lock()
    defer mutex.Unlock()

    // Find the server with the least number of active connections
    minConnections := atomic.LoadInt32(&servers[0].ActiveConnections)
    minIndex := 0
    for i, server := range servers {
        if connections := atomic.LoadInt32(&server.ActiveConnections); connections < minConnections {
            minConnections = connections
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

    // Decrement the number of active connections for the server
    atomic.AddInt32(&server.ActiveConnections, -1)
}

func main() {
    http.HandleFunc("/", handleRequest)
    http.ListenAndServe(":8000", nil)
}