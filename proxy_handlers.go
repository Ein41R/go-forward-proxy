package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

func handleGet(w http.ResponseWriter, r *http.Request) {
	handleAny(w, r)
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

func handleAny(w http.ResponseWriter, r *http.Request) {
	// strip per hop headers

	for _, h := range perHopHeaders {
		r.Header.Del(h)
	}

	headers := make(map[string]string)
	for key, values := range r.Header {
		headers[key] = values[0]
	}

	response := MakeRequest(r.URL.String(), r.Method, headers)
	// log.Printf("Response from %s: %s\n", r.URL.String(), response)
	w.Write([]byte(response))
}

func MakeRequest(URL string, method string, headers map[string]string) string {
	log.Printf("Making %s request\n", method)
	client := &http.Client{}
	req, _ := http.NewRequest(method, URL, nil)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("Err is", err)
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)
	response := string(resBody)

	log.Println("Request successfuly made to", URL)
	return response
}
