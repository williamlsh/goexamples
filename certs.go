package x // import github.com/williamlsh/x

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	organization = "Company, INC."
	validFrom    = "Jun 7 15:04:05 2021"
	validFor     = 365 * 24 * time.Hour
	host         = "127.0.0.1"

	caFile   = "ca.pem"
	certFile = "cert.pem"
	keyFile  = "key.pem"
)

// GenerateCerts generates self signed CA certificate and server side key pair.
func GenerateCerts() error {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %v", err)
	}

	notBefore, err := time.Parse("Jan 2 15:04:05 2006", validFrom)
	if err != nil {
		return fmt.Errorf("failed to parse creation date: %w", err)
	}
	notAfter := notBefore.Add(validFor)

	// Set up CA certificate.
	ca := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{organization},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Create private and public key.
	caPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	// Create the CA.
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	caOut, err := os.Create(caFile)
	if err != nil {
		return fmt.Errorf("failed to open ca.pem for writing: %w", err)
	}
	defer caOut.Close()

	// pem encode.
	err = pem.Encode(caOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return fmt.Errorf("error pem encode: %w", err)
	}
	fmt.Println("Wrote ca.pem")

	// Set up our server certificate.
	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{organization},
		},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("invalid host: %s", host)
	}
	cert.IPAddresses = append(cert.IPAddresses, ip, net.IPv6loopback)

	certPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return fmt.Errorf("failed to open cert.pem for writing: %w", err)
	}
	defer certOut.Close()

	err = pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return fmt.Errorf("error pem encode: %w", err)
	}
	fmt.Println("Wrote cert.pem")

	privBytes, err := x509.MarshalPKCS8PrivateKey(certPrivKey)
	if err != nil {
		return fmt.Errorf("unable to marshal private key: %w", err)
	}

	keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key.pem for writing: %w", err)
	}
	defer keyOut.Close()

	err = pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})
	if err != nil {
		return fmt.Errorf("error pem encode: %w", err)
	}
	fmt.Println("Wrote key.pem")

	return nil
}

func removeCerts() error {
	for _, file := range []string{
		caFile,
		certFile,
		keyFile,
	} {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	return nil
}
