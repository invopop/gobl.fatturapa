package sdi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/http2"
)

// ServerConfig defines soap server configuration options
type ServerConfig struct {
	Host       string
	Port       string
	Verbose    bool
	CACert     *x509.CertPool
	CertAuth   tls.Certificate
	ClientAuth tls.ClientAuthType
}

// RunServer sets up a server for receiving invoices from SdI
func RunServer(config *ServerConfig, messageHandler http.HandlerFunc) error {
	tlsConfig := &tls.Config{
		ClientAuth:   config.ClientAuth,
		ClientCAs:    config.CACert,
		Certificates: []tls.Certificate{config.CertAuth},
	}

	server := http.Server{
		Addr:      config.Host + ":" + config.Port,
		Handler:   messageHandler,
		TLSConfig: tlsConfig,
	}

	if err := http2.ConfigureServer(&server, &http2.Server{}); err != nil {
		log.Fatalf("Failed to configure HTTP/2: %v", err)
	}

	go func() {
		if err := server.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")

	return nil
}
