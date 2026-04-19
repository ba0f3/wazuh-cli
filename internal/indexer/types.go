package indexer

import "encoding/json"

// SearchResponse represents the response from an OpenSearch /_search request.
type SearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits         Hits                   `json:"hits"`
	Aggregations map[string]Aggregation `json:"aggregations,omitempty"`
}

// Hits represents the wrapper around the actual hit documents.
type Hits struct {
	Total struct {
		Value    int64  `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

// Hit represents a single document from OpenSearch.
type Hit struct {
	Index  string          `json:"_index"`
	ID     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

// Aggregation represents an aggregation result bucket.
type Aggregation struct {
	Buckets []Bucket `json:"buckets"`
}

// Bucket represents a single bucket in an aggregation.
type Bucket struct {
	Key      interface{} `json:"key"`
	DocCount int64       `json:"doc_count"`
}

// CountResponse represents the response from an OpenSearch /_count request.
type CountResponse struct {
	Count  int64 `json:"count"`
	Shards struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

// ErrorResponse represents an OpenSearch error response.
type ErrorResponse struct {
	Error struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error"`
	Status int `json:"status"`
}
