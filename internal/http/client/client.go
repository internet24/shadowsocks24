package client

import (
	"crypto/tls"
	"github.com/internet24/shadowsocks24/internal/config"
	"net/http"
	"time"
)

// New creates an instance of http.Client with custom configuration.
func New(c *config.Config) *http.Client {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return &http.Client{
		Transport: customTransport,
		Timeout:   time.Duration(c.HttpClient.Timeout) * time.Second,
	}
}
