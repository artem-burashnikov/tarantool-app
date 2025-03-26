// Structures used internally to handle http requests and responses.

package domain

import (
	"encoding/json"
)

// POST Request.
type AppPostRequest struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

// POST Response.
type AppPostResponse struct {
	Message string
	Key     string
	Size    uint
}

// GET Response.
type AppGetResponse struct {
	Value json.RawMessage `json:"value"`
}
