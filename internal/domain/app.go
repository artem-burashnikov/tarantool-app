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
	Message string `json:"message"`
	Key     string `json:"key"`
	Size    uint   `json:"size"`
}

// GET Response.
type AppGetResponse struct {
	Value json.RawMessage `json:"value"`
}

// PUT Response
type AppPutResponse struct {
	Message string `json:"message"`
	Key     string `json:"key"`
	Size    uint   `json:"size"`
}

// DELETE Response
type AppDeleteResponse struct {
	Message string `json:"message"`
	Key     string `json:"key"`
}
