package main

import (
	// "io"

	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
	// "https://github.com/pmezard/adblock"
	// "os"
)

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
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleConnect(w http.ResponseWriter, r *http.Request) { //test using curl -x http://127.0.0.1:1234 https://google.com
	var wg sync.WaitGroup

	// 1-2 	client connects to proxy over tcp, need to hijcak connection to handle CONNECT method
	client_conn, bufrw, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bufrw.Flush()             // flush any buffered data to client
	defer client_conn.Close() // close connection once function exits
	// 3 	client sends CONNECT <host> HTTP/1.1
	log.Printf("Outgoing CONNECT request for %s\n", r.Host)

	// todo: 4 proxy checks Access control list and blocks connection if neccessary

	// 5 	proxy connects to host via tcp
	host_conn, err := net.DialTimeout("tcp", r.Host, 3*time.Second) //timeouts after 3 seconds
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer host_conn.Close() // close connection once function exits
	log.Printf("TCP connection to %s established\n", r.Host)

	// 6 	proxy responds with HTTP/1.1 200 Connection established
	_, err = client_conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		return
	}

	log.Printf("%v -> %v", client_conn.LocalAddr(), host_conn.RemoteAddr())

	// 7 	proxy enters pipe mode
	wg.Go(func() { pipe(client_conn, host_conn) })
	wg.Go(func() { pipe(host_conn, client_conn) })
	wg.Wait() // wait for both goroutines to finish
}

func pipe(src io.Writer, dst io.Reader) {
	written, err := io.Copy(src, dst)
	log.Printf("Piped %d bytes", written)
	if err != nil {
		log.Printf("Error occurred while piping data: %v", err)
	}
}
