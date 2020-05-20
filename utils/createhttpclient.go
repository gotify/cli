package utils

import (
	"crypto/tls"
	"net/http"
	"os"
	"strings"
)

func CreateHTTPClient() *http.Client {
	skipVerify := strings.ToLower(os.Getenv("GOTIFY_SKIP_VERIFY_TLS")) == "true"
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: skipVerify}
	return &http.Client{Transport: customTransport}
}
