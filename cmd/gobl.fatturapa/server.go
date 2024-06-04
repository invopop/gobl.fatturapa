package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"

	sdi "github.com/invopop/gobl.fatturapa/sdi"
	"github.com/spf13/cobra"
)

type serverOpts struct {
	*rootOpts
}

func server(o *rootOpts) *serverOpts {
	return &serverOpts{rootOpts: o}
}

func (c *serverOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server [host] [port]",
		Short: "Server for communication with SdI in selected environment",
		RunE:  c.runE,
	}

	f := cmd.Flags()
	f.BoolP("verbose", "v", false, "Logs all requests into the console")
	f.String("ca-cert", "", "Path to a file containing the CA certificate")
	f.String("cert", "", "Path to a file containing the certificate")
	f.String("key", "", "Path to a file containing the certificate key")
	f.String("client-auth", "RequireAndVerifyClientCert", "Client authentication type (NoClientCert, RequestClientCert, RequireAnyClientCert, VerifyClientCertIfGiven, RequireAndVerifyClientCert)")
	_ = cmd.MarkFlagRequired("ca-cert")
	_ = cmd.MarkFlagRequired("cert")
	_ = cmd.MarkFlagRequired("key")

	return cmd
}

func (c *serverOpts) runE(cmd *cobra.Command, args []string) error {
	host := inputHost(args)
	port := inputPort(args)

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	caCertPEM, err := loadDataFromFlag(cmd, "ca-cert")
	if err != nil {
		return err
	}

	certPEM, err := loadDataFromFlag(cmd, "cert")
	if err != nil {
		return err
	}

	keyPEM, err := loadDataFromFlag(cmd, "key")
	if err != nil {
		return err
	}

	tlsCert, err := loadTLSCert(certPEM, keyPEM)
	if err != nil {
		return err
	}

	caCertPool, err := loadCertPoolFromPEM(caCertPEM)
	if err != nil {
		return err
	}

	clientAuth, err := getClientAuthTypeFromFlag(cmd)
	if err != nil {
		return err
	}

	if verbose {
		log.Printf("Server start: %s:%s\n", host, port)
		log.Printf("Client auth: %s\n", clientAuth)
	}

	config := &sdi.ServerConfig{
		Host:       host,
		Port:       port,
		Verbose:    verbose,
		CACert:     caCertPool,
		CertAuth:   tlsCert,
		ClientAuth: clientAuth,
	}

	err = sdi.RunServer(config, sdi.MessageHandler)
	if err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func inputHost(args []string) string {
	if len(args) > 0 && args[0] != "-" {
		return args[0]
	}
	return ""
}

func inputPort(args []string) string {
	if len(args) > 1 && args[1] != "-" {
		return args[1]
	}
	return ""
}

func loadTLSCert(publicCertPEM, privateKeyPEM []byte) (tls.Certificate, error) {
	serverTLSCert, err := tls.X509KeyPair(publicCertPEM, privateKeyPEM)
	if err != nil {
		return serverTLSCert, err
	}

	return serverTLSCert, nil
}

func loadCertPoolFromPEM(caCertPEM []byte) (*x509.CertPool, error) {
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		caCertPool = x509.NewCertPool()
	}

	if !caCertPool.AppendCertsFromPEM(caCertPEM) {
		return nil, fmt.Errorf("no certificates appended")
	}

	return caCertPool, nil
}

var clientAuthTypes = map[string]tls.ClientAuthType{
	"NoClientCert":               tls.NoClientCert,
	"RequestClientCert":          tls.RequestClientCert,
	"RequireAnyClientCert":       tls.RequireAnyClientCert,
	"VerifyClientCertIfGiven":    tls.VerifyClientCertIfGiven,
	"RequireAndVerifyClientCert": tls.RequireAndVerifyClientCert,
}

func getClientAuthTypeFromFlag(cmd *cobra.Command) (tls.ClientAuthType, error) {
	flagName := "client-auth"
	clientAuthTypeName, err := cmd.Flags().GetString(flagName)
	if err != nil {
		return tls.NoClientCert, fmt.Errorf("failed to get value of flag %s: %w", flagName, err)
	}

	clientAuthType, ok := clientAuthTypes[clientAuthTypeName]
	if !ok {
		return tls.NoClientCert, fmt.Errorf("invalid client authentication type: %s", clientAuthTypeName)
	}

	return clientAuthType, nil
}
