package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// QueryParams contains common Wazuh API query parameters.
type QueryParams struct {
	Offset  int
	Limit   int
	Sort    string
	Search  string
	Query   string // WQL
	Select  []string
	Filters map[string]string
}

// Get performs an authenticated GET request and returns the parsed response.
func (c *Client) Get(path string, params *QueryParams) (*Response, error) {
	return c.do(http.MethodGet, path, params, nil)
}

// Post performs an authenticated POST request with a JSON body.
func (c *Client) Post(path string, body interface{}) (*Response, error) {
	return c.do(http.MethodPost, path, nil, body)
}

// Put performs an authenticated PUT request with a JSON body.
func (c *Client) Put(path string, body interface{}) (*Response, error) {
	return c.do(http.MethodPut, path, nil, body)
}

// Delete performs an authenticated DELETE request.
func (c *Client) Delete(path string, params *QueryParams) (*Response, error) {
	return c.do(http.MethodDelete, path, params, nil)
}

// do executes the HTTP request with auth, debug logging, and error handling.
func (c *Client) do(method, path string, params *QueryParams, body interface{}) (*Response, error) {
	// Build full URL with query parameters
	fullURL := c.baseURL + path
	if params != nil {
		q, err := buildQuery(params)
		if err != nil {
			return nil, err
		}
		if len(q) > 0 {
			fullURL += "?" + q.Encode()
		}
	}

	// Build request body
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Attach auth token
	token, err := c.auth.Token()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	if c.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s %s\n", method, fullURL)
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
		fmt.Fprintf(os.Stderr, "[DEBUG] Response %d: %s\n", resp.StatusCode, truncate(string(respBody), 500))
	}

	// Handle 401 — token may have expired mid-session
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, &AuthError{Message: "token expired or invalid (try again to re-authenticate)"}
	}

	return parseResponse(respBody, resp.StatusCode)
}

// buildQuery converts QueryParams into url.Values.
func buildQuery(p *QueryParams) (url.Values, error) {
	q := url.Values{}
	if p.Offset > 0 {
		q.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Sort != "" {
		q.Set("sort", p.Sort)
	}
	if p.Search != "" {
		q.Set("search", p.Search)
	}
	if p.Query != "" {
		q.Set("q", p.Query)
	}
	for i, s := range p.Select {
		if i == 0 {
			q.Set("select", s)
		} else {
			q.Add("select", s)
		}
	}
	for k, v := range p.Filters {
		q.Set(k, v)
	}
	return q, nil
}

// truncate shortens a string for debug logging.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
