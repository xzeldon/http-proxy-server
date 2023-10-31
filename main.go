package main

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	PORT          string
	AUTH_USERNAME string
	AUTH_PASSWORD string
)

func checkProxyAuth(r *http.Request, username, password string) bool {
	// If no username or password provided, always allow
	if AUTH_USERNAME == "" && AUTH_PASSWORD == "" {
		return true
	}
	authHeader := r.Header.Get("Proxy-Authorization")
	const prefix = "Basic "
	if !strings.HasPrefix(authHeader, prefix) {
		return false
	}

	payload, err := base64.StdEncoding.DecodeString(authHeader[len(prefix):])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return false
	}

	return pair[0] == AUTH_USERNAME && pair[1] == AUTH_PASSWORD
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if !checkProxyAuth(r, AUTH_USERNAME, AUTH_PASSWORD) {
		w.Header().Set("Proxy-Authenticate", `Basic realm="Provide username and password"`)
		http.Error(w, "Proxy authentication required", http.StatusProxyAuthRequired)
		return
	}

	if r.Method == http.MethodConnect {
		handleTunneling(w, r)
	} else {
		handleHTTP(w, r)
	}
}

func logRequest(r *http.Request) {
	log.Printf("Method: %s, URL: %s", r.Method, r.URL.String())
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	logRequest(req)

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func main() {
	flag.StringVar(&PORT, "port", "3000", "Specify the port the proxy will run on")
	flag.StringVar(&AUTH_USERNAME, "username", "", "Username for proxy authentication")
	flag.StringVar(&AUTH_PASSWORD, "password", "", "Password for proxy authentication")
	flag.Parse()

	if (AUTH_USERNAME == "" && AUTH_PASSWORD != "") || (AUTH_USERNAME != "" && AUTH_PASSWORD == "") {
		log.Fatal("Error: Both username and password must be provided, or neither should be.")
	}

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: http.HandlerFunc(mainHandler),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Printf("Starting proxy server on :%s\n", PORT)
	log.Fatal(server.ListenAndServe())
}
