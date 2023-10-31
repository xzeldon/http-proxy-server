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

type ProxyServer struct {
	port     string
	username string
	password string
}

func (ps *ProxyServer) isAuthRequired() bool {
	return ps.username != "" && ps.password != ""
}

func (ps *ProxyServer) checkProxyAuth(r *http.Request) bool {
	if !ps.isAuthRequired() {
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

	return pair[0] == ps.username && pair[1] == ps.password
}

func (ps *ProxyServer) requestHandler(w http.ResponseWriter, r *http.Request) {
	if !ps.checkProxyAuth(r) {
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

// Manages tunneling requests (used for HTTPS connections).
func handleTunneling(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)

	// Hijack the connection to manage it at the TCP level
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Transfer data between client and destination
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
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

func logRequest(r *http.Request) {
	log.Printf("Method: %s, URL: %s", r.Method, r.URL.String())
}

func main() {
	var port, username, password string

	flag.StringVar(&port, "port", "3000", "Specify the port the proxy will run on")
	flag.StringVar(&username, "username", "", "Username for proxy authentication")
	flag.StringVar(&password, "password", "", "Password for proxy authentication")
	flag.Parse()

	if (username == "" && password != "") || (username != "" && password == "") {
		log.Fatal("Error: Both username and password must be provided, or neither should be.")
	}

	proxy := &ProxyServer{port: port, username: username, password: password}
	server := &http.Server{
		Addr:    ":" + proxy.port,
		Handler: http.HandlerFunc(proxy.requestHandler),
		// Disable HTTP/2 (https://github.com/golang/go/issues/14797#issuecomment-196103814)
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Printf("Starting proxy server on :%s\n", proxy.port)
	log.Fatal(server.ListenAndServe())
}
