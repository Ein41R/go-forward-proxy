package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// WARNING: incomplete list of hop-by-hop headers
// which will be stripped
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
	ctx := context.Background()
	ctx, err := loadConfig(ctx)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	//EXPLINATION: ctx.Value returns interface, which we assert the Config type
	config, ok := ctx.Value(cfgInterfaceKey).(Config)
	if !ok {
		log.Fatal("Failed to load config")
	}

	host := config.Host
	port := config.Port

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
