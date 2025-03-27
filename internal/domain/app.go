// Structures used internally to handle http requests and responses.

package domain

import (
	"encoding/json"
)

// POST Request from user.
type AppPostRequest struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

// POST Response for user.
type AppPostResponse struct {
	Message string `json:"message"`
	Key     string `json:"key"`
	Size    uint   `json:"size"`
}

// GET Request from user.
type AppGetRequest struct {
	Key string `json:"key"`
}

// GET Response for user.
type AppGetResponse struct {
	Value json.RawMessage `json:"value"`
}

// PUT Request from user.
type AppPutRequest struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

// PUT Response for user.
type AppPutResponse struct {
	Message string `json:"message"`
	Key     string `json:"key"`
	Size    uint   `json:"size"`
}

// DELETE Response
type AppDeleteRequest struct {
	Key string `json:"key"`
}

// DELETE Response
type AppDeleteResponse struct {
	Message string `json:"message"`
	Key     string `json:"key"`
}
