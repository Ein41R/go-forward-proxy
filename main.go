package main

import (
	// "io"
	"log"
	"net/http"
	// "https://github.com/pmezard/adblock"
	// "os"
)

func main() {
	host := "127.0.0.1"
	port := "1234"
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleFunc)

	log.Printf("Server started at %s:%s\n", host, port)
	err := http.ListenAndServe(host+":"+port, mux)
	if err != nil {
		panic(err)
	}
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
	case http.MethodConnect:
		handleConnect(w, r)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleConnect(w http.ResponseWriter, r *http.Request) {

}
