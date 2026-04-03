package main

import (
	// "io"

	"log"
	"net/http"
	// "https://github.com/pmezard/adblock"
	// "os"
)

// TODO: parse from config.json
var perHopHeaders = []string{
	"Proxy-Connection",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Connection",
	"Keep-Alive",
	"TE",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

func main() {
	host := "127.0.0.1"
	port := "1234"
	// mux := http.NewServeMux()

	// mux.HandleFunc("/", handleFunc)

	log.Printf("Server started at %s:%s\n", host, port)
	err := http.ListenAndServe(host+":"+port, http.HandlerFunc(handleFunc))
	if err != nil {
		panic(err)
	}
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodConnect:
		handleConnect(w, r)
	default:
		handleAny(w, r)
	}
}
