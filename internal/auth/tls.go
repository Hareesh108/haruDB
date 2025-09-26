// internal/auth/tls.go
package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// TLSManager handles TLS configuration and certificate management
type TLSManager struct {
	certFile string
	keyFile  string
	config   *tls.Config
}

// NewTLSManager creates a new TLS manager
func NewTLSManager(dataDir string) *TLSManager {
	certFile := filepath.Join(dataDir, "server.crt")
	keyFile := filepath.Join(dataDir, "server.key")

	tm := &TLSManager{
		certFile: certFile,
		keyFile:  keyFile,
	}

	// Generate self-signed certificate if it doesn't exist
	if !tm.certificateExists() {
		if err := tm.generateSelfSignedCert(); err != nil {
			fmt.Printf("Warning: Failed to generate self-signed certificate: %v\n", err)
		}
	}

	// Load TLS configuration
	tm.loadTLSConfig()

	return tm
}

// GetTLSConfig returns the TLS configuration
func (tm *TLSManager) GetTLSConfig() *tls.Config {
	return tm.config
}

// certificateExists checks if certificate files exist
func (tm *TLSManager) certificateExists() bool {
	_, certErr := os.Stat(tm.certFile)
	_, keyErr := os.Stat(tm.keyFile)
	return certErr == nil && keyErr == nil
}

// generateSelfSignedCert generates a self-signed certificate
func (tm *TLSManager) generateSelfSignedCert() error {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"HaruDB"},
			Country:       []string{"IN"},
			Province:      []string{""},
			Locality:      []string{"Pune"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:    []string{"localhost"},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Save certificate file
	certOut, err := os.Create(tm.certFile)
	if err != nil {
		return fmt.Errorf("failed to open cert file for writing: %w", err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	// Save private key file
	keyOut, err := os.OpenFile(tm.keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key file for writing: %w", err)
	}
	defer keyOut.Close()

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	return nil
}

// loadTLSConfig loads the TLS configuration
func (tm *TLSManager) loadTLSConfig() {
	cert, err := tls.LoadX509KeyPair(tm.certFile, tm.keyFile)
	if err != nil {
		fmt.Printf("Warning: Failed to load TLS certificate: %v\n", err)
		tm.config = nil
		return
	}

	tm.config = &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

// IsTLSEnabled returns true if TLS is properly configured
func (tm *TLSManager) IsTLSEnabled() bool {
	return tm.config != nil
}
