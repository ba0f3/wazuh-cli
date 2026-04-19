package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ba0f3/wazuh-cli/internal/config"
)

// Client wraps an http.Client with Wazuh-specific auth and configuration.
type Client struct {
	cfg        *config.Config
	httpClient *http.Client
	auth       *authManager
	baseURL    string
	debug      bool
}

// NewClient creates a new Wazuh API client from the given config.
func NewClient(cfg *config.Config) (*Client, error) {
	tlsConfig, err := buildTLSConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("building TLS config: %w", err)
	}

	if cfg.Insecure {
		fmt.Fprintln(os.Stderr, "WARNING: TLS certificate verification is disabled (--insecure)")
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
		Transport: transport,
	}

	// Normalise base URL (strip trailing slash)
	baseURL := cfg.URL
	for len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	c := &Client{
		cfg:        cfg,
		httpClient: httpClient,
		baseURL:    baseURL,
		debug:      cfg.Debug,
	}

	c.auth = newAuthManager(c, cfg)
	return c, nil
}

// buildTLSConfig creates a tls.Config based on the client configuration.
func buildTLSConfig(cfg *config.Config) (*tls.Config, error) {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: cfg.Insecure, //nolint:gosec // intentional, flagged by --insecure
	}

	// Custom CA certificate
	if cfg.CACert != "" {
		caCert, err := os.ReadFile(cfg.CACert)
		if err != nil {
			return nil, fmt.Errorf("reading CA cert %s: %w", cfg.CACert, err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("parsing CA cert %s: no valid certificates found", cfg.CACert)
		}
		tlsCfg.RootCAs = pool
	}

	// Client certificate (mTLS)
	if cfg.ClientCert != "" && cfg.ClientKey != "" {
		cert, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("loading client cert/key: %w", err)
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
	}

	return tlsCfg, nil
}

// BaseURL returns the configured API base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// HTTPClient returns the underlying http.Client (for testing).
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// Auth returns the authentication manager for this client.
func (c *Client) Auth() *authManager {
	return c.auth
}
