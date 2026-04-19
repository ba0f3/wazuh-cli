package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ba0f3/wazuh-cli/internal/config"
)

// tokenCache is what gets stored on disk.
type tokenCache struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// authManager handles JWT acquisition and token caching.
type authManager struct {
	client    *Client
	cfg       *config.Config
	tokenPath string
	cached    *tokenCache
}

// newAuthManager creates a new authManager.
func newAuthManager(c *Client, cfg *config.Config) *authManager {
	return &authManager{
		client:    c,
		cfg:       cfg,
		tokenPath: config.DefaultTokenPath(),
	}
}

// Token returns a valid JWT, refreshing or re-authenticating as needed.
func (a *authManager) Token() (string, error) {
	// 1. Use raw token if provided directly (--raw-token or WAZUH_TOKEN)
	if a.cfg.RawToken != "" {
		return a.cfg.RawToken, nil
	}

	// 2. Use cached token if still valid (with 30s buffer)
	if a.cached != nil && time.Now().Add(30*time.Second).Before(a.cached.ExpiresAt) {
		return a.cached.Token, nil
	}

	// 3. Try loading from disk cache
	if disk, err := a.loadTokenCache(); err == nil {
		if time.Now().Add(30 * time.Second).Before(disk.ExpiresAt) {
			a.cached = disk
			return disk.Token, nil
		}
	}

	// 4. Authenticate fresh
	return a.authenticate()
}

// authenticate fetches a new JWT from the Wazuh API.
func (a *authManager) authenticate() (string, error) {
	token, err := a.authenticateBasic("GET")
	if err != nil {
		// Fallback to POST if GET fails
		token, err = a.authenticateBasic("POST")
	}
	if err != nil {
		// Fallback to JSON body POST
		token, err = a.authenticateJSON()
	}
	return token, err
}

func (a *authManager) authenticateBasic(method string) (string, error) {
	url := a.client.baseURL + "/security/user/authenticate"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("building auth request: %w", err)
	}

	req.SetBasicAuth(a.cfg.User, a.cfg.Password)

	if a.client.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s %s (Basic Auth: %s:***, pass_len: %d)\n", method, url, a.cfg.User, len(a.cfg.Password))
	}

	return a.doAuthRequest(req)
}

func (a *authManager) authenticateJSON() (string, error) {
	url := a.client.baseURL + "/security/user/authenticate"
	body, _ := json.Marshal(map[string]string{
		"user":     a.cfg.User,
		"password": a.cfg.Password,
	})

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	if a.client.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] POST %s (JSON Body, user: %s)\n", url, a.cfg.User)
	}

	return a.doAuthRequest(req)
}

func (a *authManager) doAuthRequest(req *http.Request) (string, error) {
	resp, err := a.client.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response envelope
	var envelope struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
		Error   int    `json:"error"`
		Message string `json:"message"`
		Title   string `json:"title"`
		Detail  string `json:"detail"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", &AuthError{Message: "invalid credentials"}
		}
		return "", fmt.Errorf("authentication failed with status %d", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		msg := envelope.Detail
		if msg == "" {
			msg = envelope.Message
		}
		if msg == "" {
			msg = envelope.Title
		}
		if msg == "" {
			msg = fmt.Sprintf("status %d", resp.StatusCode)
		}
		return "", &AuthError{Message: msg}
	}

	if envelope.Error != 0 {
		return "", &AuthError{Message: envelope.Message}
	}
	if envelope.Data.Token == "" {
		return "", &AuthError{Message: "empty token in response"}
	}

	// Cache for 15 minutes
	tc := &tokenCache{
		Token:     envelope.Data.Token,
		ExpiresAt: time.Now().Add(14 * time.Minute),
	}
	a.cached = tc
	_ = a.saveTokenCache(tc)

	if a.client.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Authenticated successfully, token cached until %s\n", tc.ExpiresAt.Format(time.RFC3339))
	}

	return tc.Token, nil
}

// Logout invalidates the current session token.
func (a *authManager) Logout() error {
	url := a.client.baseURL + "/security/user/authenticate"
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	token, err := a.Token()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := a.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Clear caches
	a.cached = nil
	_ = os.Remove(a.tokenPath)
	return nil
}

// loadTokenCache reads the token cache from disk.
func (a *authManager) loadTokenCache() (*tokenCache, error) {
	data, err := os.ReadFile(a.tokenPath)
	if err != nil {
		return nil, err
	}
	var tc tokenCache
	if err := json.Unmarshal(data, &tc); err != nil {
		return nil, err
	}
	return &tc, nil
}

// saveTokenCache writes the token cache to disk with 0600 permissions.
func (a *authManager) saveTokenCache(tc *tokenCache) error {
	if err := os.MkdirAll(filepath.Dir(a.tokenPath), 0700); err != nil {
		return err
	}
	data, err := json.Marshal(tc)
	if err != nil {
		return err
	}
	return os.WriteFile(a.tokenPath, data, 0600)
}

// TokenPath returns the path to the cached token file.
func (a *authManager) TokenPath() string {
	return a.tokenPath
}

// AuthError represents an authentication failure.
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return "authentication failed: " + e.Message
}
