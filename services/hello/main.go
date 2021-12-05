package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request from: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "Hello you:! %v", r.URL.Path[1:])
}
