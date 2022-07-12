package main

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
)

func ListenAndServe(listenAddress, certFile, keyFile string, handler http.HandlerFunc) error {
	// Load certs
	var err error
	certs := make([]tls.Certificate, 1)
	certs[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	// We currently only use the cert-related stuff from tls.Config,
	// so we don't need to make a full copy.
	config := &tls.Config{
		Certificates: certs,
	}

	// Open the listeners
	udpAddr, err := net.ResolveUDPAddr("udp", listenAddress)
	if err != nil {
		return err
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	defer udpConn.Close()

	// Start the servers
	httpServer := &http.Server{
		Addr:      listenAddress,
		TLSConfig: config,
	}

	quicServer := &http3.Server{
		Server: httpServer,
	}

	httpServer.Handler = handler

	hErr := make(chan error)
	qErr := make(chan error)
	go func() {
		qErr <- quicServer.Serve(udpConn)
	}()

	select {
	case err := <-hErr:
		quicServer.Close()
		return err
	case err := <-qErr:
		// Cannot close the HTTP server or wait for requests to complete properly :/
		return err
	}
}
