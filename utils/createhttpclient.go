package utils

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"strings"
)

func CreateHTTPClient() *http.Client {
	skipVerify := strings.ToLower(os.Getenv("GOTIFY_SKIP_VERIFY_TLS")) == "true"
	certFile := os.Getenv("SSL_CERT_FILE")
	if skipVerify && certFile != "" {
		Exit1With("GOTIFY_SKIP_VERIFY_TLS and SSL_CERT_FILE shouldn't be set at the same time")
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	rootCAs := customTransport.TLSClientConfig.RootCAs
	if certFile != "" {
		cert, err := os.ReadFile(certFile)
		if err != nil {
			Exit1With("Failed to read cert:", err)
		}
		rootCAs = x509.NewCertPool()
		ok := rootCAs.AppendCertsFromPEM(cert)
		if !ok {
			Exit1With("Failed to parse cert", certFile)
		}
	}
	customTransport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: skipVerify,
		RootCAs: rootCAs,
	}
	return &http.Client{Transport: customTransport}
}
