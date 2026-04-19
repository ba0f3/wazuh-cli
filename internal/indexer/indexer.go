package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/ba0f3/wazuh-cli/internal/config"
)

// Client is a lightweight OpenSearch client for querying Wazuh alerts.
type Client struct {
	cfg        *config.Config
	httpClient *http.Client
	baseURL    string
	username   string
	password   string
	debug      bool
}

// NewClient creates a new OpenSearch client using the indexer configuration.
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.IndexerURL == "" {
		return nil, fmt.Errorf("indexer_url is not configured")
	}

	tlsConfig, err := client.BuildTLSConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("building TLS config for indexer: %w", err)
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
		Transport: transport,
	}

	// Normalise base URL (strip trailing slash)
	baseURL := cfg.IndexerURL
	for len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	user := cfg.IndexerUser
	if user == "" {
		user = cfg.User // fallback to Manager user
	}
	pass := cfg.IndexerPassword
	if pass == "" {
		pass = cfg.Password // fallback to Manager password
	}

	return &Client{
		cfg:        cfg,
		httpClient: httpClient,
		baseURL:    baseURL,
		username:   user,
		password:   pass,
		debug:      cfg.Debug,
	}, nil
}

// Search executes an OpenSearch _search query against the given index.
func (c *Client) Search(index string, queryDSL interface{}) (*SearchResponse, error) {
	path := fmt.Sprintf("/%s/_search", index)

	respBody, err := c.do(http.MethodPost, path, queryDSL)
	if err != nil {
		return nil, err
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(respBody, &searchResp); err != nil {
		return nil, fmt.Errorf("parsing search response: %w", err)
	}

	return &searchResp, nil
}

// Get executes an OpenSearch _doc GET request by ID.
func (c *Client) Get(index string, id string) (*Hit, error) {
	path := fmt.Sprintf("/%s/_doc/%s", index, id)

	respBody, err := c.do(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var hit Hit
	if err := json.Unmarshal(respBody, &hit); err != nil {
		return nil, fmt.Errorf("parsing get response: %w", err)
	}

	if hit.ID == "" {
		return nil, fmt.Errorf("document not found")
	}

	return &hit, nil
}

// Count executes an OpenSearch _count request.
func (c *Client) Count(index string, queryDSL interface{}) (*CountResponse, error) {
	path := fmt.Sprintf("/%s/_count", index)

	respBody, err := c.do(http.MethodPost, path, queryDSL)
	if err != nil {
		return nil, err
	}

	var countResp CountResponse
	if err := json.Unmarshal(respBody, &countResp); err != nil {
		return nil, fmt.Errorf("parsing count response: %w", err)
	}

	return &countResp, nil
}

// do executes the HTTP request with Basic Auth.
func (c *Client) do(method, path string, body interface{}) ([]byte, error) {
	fullURL := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)

		if c.debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Request Body: %s\n", string(data))
		}
	}

	req, err := http.NewRequestWithContext(context.Background(), method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	if c.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] INDEXER %s %s\n", method, fullURL)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if c.debug {
		truncateMax := 500
		respStr := string(respBody)
		if len(respStr) > truncateMax {
			respStr = respStr[:truncateMax] + "..."
		}
		fmt.Fprintf(os.Stderr, "[DEBUG] INDEXER Response %d: %s\n", resp.StatusCode, respStr)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if jsonErr := json.Unmarshal(respBody, &errResp); jsonErr == nil && errResp.Error.Type != "" {
			return nil, fmt.Errorf("indexer error (%d): %s - %s", resp.StatusCode, errResp.Error.Type, errResp.Error.Reason)
		}

		// If it's a 404 on Get, that's handled cleanly by returning empty struct or checking ID.
		if method == http.MethodGet && resp.StatusCode == 404 && strings.Contains(path, "/_doc/") {
			return []byte("{}"), nil // Return empty JSON, letting caller handle it
		}

		return nil, fmt.Errorf("indexer API error: status %d", resp.StatusCode)
	}

	return respBody, nil
}
