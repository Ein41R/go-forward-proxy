package main

import (
	// "io"
	"context"
	"fmt"
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
	//context with config values
	ctx := context.Background()
	ctx, err := loadConfig(ctx)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	config, ok := ctx.Value(cfgInterfaceKey).(Config)
	//type assertion since ctx.Value returns interface{}
	//using private typed key for typesafety
	if !ok {
		log.Fatal("Failed to load config")
	}

	host := config.Host
	port := config.Port
	// mux := http.NewServeMux()

	// mux.HandleFunc("/", handleFunc)

	log.Printf("Server started at %s:%d\n", host, port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), http.HandlerFunc(handleFunc))
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
