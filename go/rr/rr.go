package main

import (
    "fmt"
    "net/http"
    "net/http/httputil"
    "net/url"
    "sync"

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
    serverNames := []string{"server_1", "server_2", "server_3"}
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

func handleRequest(w http.ResponseWriter, r *http.Request) {
    server := getNextServer()
    fmt.Printf("Received connection, forwarding to server at %s\n", server.URL.String())
    r.Header.Set("X-Forwarded-For", server.URL.Hostname())
    r.Header.Set("X-Forwarded-Port", server.URL.Port())
    server.ReverseURL.ServeHTTP(w, r)
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", handleRequest).Methods(http.MethodGet)

    http.ListenAndServe(":8000", router)
}