package client

import (
	"encoding/json"
	"fmt"
)

// Response represents the parsed Wazuh API response envelope.
type Response struct {
	// Data holds the actual response payload.
	Data json.RawMessage `json:"data"`
	// Error is the Wazuh API error code (0 = success).
	Error int `json:"error"`
	// Message is set when Error != 0.
	Message string `json:"message"`

	// Pagination metadata (populated from data.total_affected_items etc.)
	TotalAffectedItems int             `json:"-"`
	AffectedItems      json.RawMessage `json:"-"`
	FailedItems        json.RawMessage `json:"-"`

	// Raw HTTP status code
	HTTPStatus int `json:"-"`
	// RawBody is the full unparsed response (for --output raw)
	RawBody []byte `json:"-"`
}

// dataEnvelope matches the Wazuh API data wrapper.
type dataEnvelope struct {
	AffectedItems      json.RawMessage `json:"affected_items"`
	TotalAffectedItems int             `json:"total_affected_items"`
	FailedItems        json.RawMessage `json:"failed_items"`
	TotalFailedItems   int             `json:"total_failed_items"`
}

// APIError represents a structured Wazuh API error.
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

// parseResponse decodes a raw Wazuh API response body.
func parseResponse(body []byte, status int) (*Response, error) {
	r := &Response{
		HTTPStatus: status,
		RawBody:    body,
	}

	if err := json.Unmarshal(body, r); err != nil {
		// Not JSON — return raw
		return r, nil
	}

	// Check for API-level error
	if r.Error != 0 {
		return r, &APIError{Code: r.Error, Message: r.Message}
	}

	// Unpack the data envelope if present
	if len(r.Data) > 0 {
		var de dataEnvelope
		if err := json.Unmarshal(r.Data, &de); err == nil {
			r.TotalAffectedItems = de.TotalAffectedItems
			r.AffectedItems = de.AffectedItems
			r.FailedItems = de.FailedItems
		}
	}

	return r, nil
}

// Items returns the affected_items array from the data envelope.
// Returns raw data if affected_items is not present.
func (r *Response) Items() json.RawMessage {
	if len(r.AffectedItems) > 0 {
		return r.AffectedItems
	}
	return r.Data
}

// IsSuccess returns true if the API call succeeded.
func (r *Response) IsSuccess() bool {
	return r.Error == 0
}
