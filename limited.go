package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/netutil"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("started")

	http.HandleFunc("/", hello)

	srv := &http.Server{
		Addr: "localhost:8000",
	}

	config := &tls.Config{
		NextProtos: []string{"http/1.1"},
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalln(err)
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatalln(err)
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)

	connectionCount := 2
	limitedListener := netutil.LimitListener(tlsListener, connectionCount)

	log.Fatalln(srv.Serve(limitedListener))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, 世界")
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
