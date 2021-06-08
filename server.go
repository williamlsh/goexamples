// Reference: https://medium.com/@shaneutt/create-sign-x509-certificates-in-golang-8ac4ae49f903

package x

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"time"
)

func x() {
	// Create a Certificate Authority

	// Creating a CA which will be used to sign all of certificates using the x509 package.
	now := time.Now()
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"CN"},
			Province:      []string{"Shenzhen"},
			Locality:      []string{"Futian"},
			StreetAddress: []string{"Exhibition center"},
			PostalCode:    []string{"94016"},
		},
		NotBefore: now,
		NotAfter:  now.AddDate(1, 0, 0),
		IsCA:      true,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Generate a private key for the CA.
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	// Create the certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		panic(err)
	}

	// PEM Encode certificate and private key for signing other certificates

	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		panic(err)
	}

	caPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		panic(err)
	}

	// --------------------------------------------------

	// Generate & Signing a Certificate.

	// Create a certificate which our CA will sign.

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"CN"},
			Province:      []string{"Shenzhen"},
			Locality:      []string{"Futian"},
			StreetAddress: []string{"Exhibition center"},
			PostalCode:    []string{"94016"},
		},
		IPAddresses: []net.IP{
			net.IPv4(127, 0, 0, 1),
			net.IPv6loopback,
		},
		NotBefore:    now,
		NotAfter:     now.AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// Create a private key for the certificate
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	// Sign the certificate with the previously created CA
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		panic(err)
	}

	// PEM Encode the certificate.
	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	// --------------------------------------------------

	// WebServer Configuration.

	// Creating a Key Pair which will be used for the server configuration.
	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		panic(err)
	}

	// And a CertPool to house our certificate for client connections.
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caPEM.Bytes())

	// Create a tls.Config which will be provided to our server.
	serverTLSConf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	// Start the httptest.Server
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "success!")
	}))
	server.TLS = serverTLSConf
	server.StartTLS()
	defer server.Close()

	// --------------------------------------------------------

	// Connecting to the Server

	clientTLSConfig := &tls.Config{
		RootCAs: certPool,
	}

	transport := &http.Transport{
		TLSClientConfig: clientTLSConfig,
	}

	client := http.Client{
		Transport: transport,
	}

	resp, err := client.Get(server.URL)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
