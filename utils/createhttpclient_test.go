package utils

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newCA(t *testing.T) ([]byte, func(domain string) (ed25519.PublicKey, ed25519.PrivateKey)) {
	caPubKey, caPrivKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	cert, err := x509.CreateCertificate(rand.Reader, &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
	}, &x509.Certificate{}, caPubKey, caPrivKey)

	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	certParsed, err := x509.ParseCertificate(cert)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert}), func(domain string) (ed25519.PublicKey, ed25519.PrivateKey) {
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatalf("failed to generate key: %v", err)
		}

		cert, err := x509.CreateCertificate(rand.Reader, &x509.Certificate{
			DNSNames:     []string{domain},
			SerialNumber: big.NewInt(2),
			NotAfter:     time.Now().Add(time.Hour),
		}, certParsed, pubKey, caPrivKey)

		if err != nil {
			t.Fatalf("failed to create certificate: %v", err)
		}

		privPEM, err := x509.MarshalPKCS8PrivateKey(privKey)
		if err != nil {
			t.Fatalf("failed to marshal private key: %v", err)
		}

		return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert}),
			pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privPEM})
	}
}

func TestCreateHTTPClient(t *testing.T) {
	caPEM, signer := newCA(t)
	wrongCAPEM, wrongSigner := newCA(t)

	certPEM, certPriv := signer("gotify.local")
	wrongDomainPEM, wrongDomainPriv := signer("gotify.invalid")
	wrongCAPEM, wrongCAPriv := wrongSigner("gotify.local")

	testTrust := func(trustCert []byte, serverPEM []byte, serverKey []byte) bool {
		serverSide, clientSide := net.Pipe()

		serverCert, err := tls.X509KeyPair(serverPEM, serverKey)
		if err != nil {
			panic(err)
		}

		tlsServer := tls.Server(serverSide, &tls.Config{
			Certificates: []tls.Certificate{
				serverCert,
			},
		})

		var certFile *os.File = nil
		if trustCert != nil {
			var err error
			certFile, err = os.CreateTemp("", "GotifyTrustCert")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			certFile.Write(trustCert)
			certFile.Close()
			os.Setenv("SSL_CERT_FILE", certFile.Name())
		}

		client := CreateHTTPClient()
		client.Transport.(*http.Transport).DialContext = func(_ context.Context, network, addr string) (net.Conn, error) {
			return clientSide, nil
		}

		os.Unsetenv("SSL_CERT_FILE")
		if certFile != nil {
			os.Remove(certFile.Name())
		}

		var failed uint32 = 0
		var unexpected error

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer serverSide.Close()
			defer wg.Done()

			if err := tlsServer.Handshake(); err == nil {
				tlsServer.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			}
			tlsServer.Close()
		}()

		go func() {
			defer clientSide.Close()
			defer wg.Done()

			if _, err := client.Get("https://gotify.local"); err != nil {
				if _, ok := errors.Unwrap(err).(*tls.CertificateVerificationError); ok {
					atomic.StoreUint32(&failed, 1)
				} else {
					unexpected = err
				}
			}
		}()

		wg.Wait()
		if unexpected != nil {
			t.Fatal(unexpected)
		}

		return atomic.LoadUint32(&failed) == 0
	}

	if !testTrust(certPEM, certPEM, certPriv) {
		t.Fatal("failed to trust valid server cert")
	}

	if !testTrust(caPEM, certPEM, certPriv) {
		t.Fatal("failed to trust valid CA")
	}

	if testTrust(caPEM, wrongCAPEM, wrongCAPriv) {
		t.Fatal("trusted invalid cert")
	}

	if testTrust(caPEM, wrongDomainPEM, wrongDomainPriv) {
		t.Fatal("trusted cert with invalid domain")
	}

	if testTrust(nil, certPEM, certPriv) {
		t.Fatal("shouldn't trust server cert")
	}
}
