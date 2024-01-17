package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		port := r.Header.Get("X-Forwarded-Port")
		fmt.Fprintf(w, "Hello from server: %s", port)
	})

	http.ListenAndServe(":"+*port, nil)
}
