package main

import (
    "flag"
    "fmt"
    "net/http"
)

func main() {
    hostname := flag.String("hostname", "0.0.0.0", "hostname to listen on")
    port := flag.String("port", "8080", "port to listen on")
    flag.Parse()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        forwardedHostname := r.Header.Get("X-Forwarded-For")
        forwardedPort := r.Header.Get("X-Forwarded-Port")

        // Print connection info to terminal
        fmt.Printf("Received connection from %s:%s\n", forwardedHostname, forwardedPort)

        fmt.Fprintf(w, "Hello from server at Hostname: %s and Port: %s", forwardedHostname, forwardedPort)
    })

    fmt.Printf("Starting server on %s:%s\n", *hostname, *port)
    http.ListenAndServe(*hostname+":"+*port, nil)
}