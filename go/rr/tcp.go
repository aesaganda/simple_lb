package main

import (
    "fmt"
    "io"
    "net"
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
    for i := 0; i < 3; i++ {
        target, _ := url.Parse(fmt.Sprintf("http://%s:%d", "0.0.0.0", 8080+i))

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

func handleTCPConnection(conn net.Conn) {
    defer conn.Close()

    server := getNextServer()
    fmt.Printf("Received TCP connection, forwarding to server at %s\n", server.URL.String())

    backendConn, err := net.Dial("tcp", server.URL.Host)
    if err != nil {
        fmt.Printf("Error connecting to backend: %v\n", err)
        return
    }
    defer backendConn.Close()

    go io.Copy(backendConn, conn)
    io.Copy(conn, backendConn)
}

func main() {
    go func() {
        router := mux.NewRouter()
        router.HandleFunc("/", handleRequest).Methods(http.MethodGet)

        http.ListenAndServe(":8000", router)
    }()

    listener, err := net.Listen("tcp", ":8001")
    if err != nil {
        panic(err)
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("Error accepting connection: %v\n", err)
            continue
        }

        go handleTCPConnection(conn)
    }
}