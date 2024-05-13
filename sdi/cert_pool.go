package sdi

import (
	"crypto/x509"
	"log"
	"os"
	"path/filepath"
)

func loadCACertFromFile(certPath string) (*x509.Certificate, error) {
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(caCert)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// PrepareCertPoolFromDir returns caCertPool from selected folder
func PrepareCertPoolFromDir(dir string) (*x509.CertPool, error) {
	caCertPool := x509.NewCertPool()

	// Find files matching the "*.pem" pattern in the directory
	matches, err := filepath.Glob(filepath.Join(dir, "*.cer"))
	if err != nil {
		return nil, err
	}

	for _, match := range matches {
		cert, err := loadCACertFromFile(match)
		if err != nil {
			log.Printf("Error reading certificate %s: %v", match, err)
			continue
		}
		caCertPool.AddCert(cert)
	}

	return caCertPool, nil
}
