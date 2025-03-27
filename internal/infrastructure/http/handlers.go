// Handlers for router.

package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"tarantool-app/internal/domain"
	"tarantool-app/internal/infrastructure/logging"
	"tarantool-app/internal/repository"
	"tarantool-app/internal/usecase"

	"github.com/gin-gonic/gin"
)

// Incapsulates available handlers.
// All of them mirror CRUD operations so are a single entity.
type RequestHandler struct {
	APIHandler *usecase.CRUD
	Logger     *logging.Logger
}

func NewRequestHandler(crudHandler *usecase.CRUD, zlog *logging.Logger) *RequestHandler {
	return &RequestHandler{APIHandler: crudHandler, Logger: zlog}
}

// GET kv/{id}
// On success responds with corresponding value stored in database.
func (rh *RequestHandler) GETHandlerFunc(c *gin.Context) {
	// Check if id is ok. If it is --- fetch the data. Then return it to user.
	key := c.Param("id")

	// Empty key is forbidden.
	// Everything else is ok?
	if key == "" {
		c.String(http.StatusBadRequest, "Missing key")
		return
	}

	// Any other key is ok. Go fetch the data from repository.
	value, err := rh.APIHandler.Read(c, key)
	// Either Tarantool somehow failed or key was not found.
	if err != nil {
		// If key was not found, respond with status code 404.
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(repository.ErrNotFound.Code, gin.H{"error": "key not found"})
		} else {
			// If further processing failed, log the error and respond with status code 500.
			rh.Logger.Warn("Tarantool failed to retreive data by key",
				"key", key,
			)
			c.String(http.StatusInternalServerError, "Internal server error.")
		}
		return
	}

	// Key was found --- respond in JSON.
	c.JSON(http.StatusOK, domain.AppGetResponse{
		Value: value},
	)
}

// POST /kv body: {key: "key", "value": {ARBITRARY JSON}}.
// On success responds with some metainfo about the inserted enrty.
func (rh *RequestHandler) POSTHandlerFunc(c *gin.Context) {
	// Verify that body is a valid JSON.
	// If so, process it further and let Tarantool create an entry.
	// Incorrect request body returns status code 400.
	var req domain.AppPostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "error: invalid request format")
		return
	}

	// Either Tarantool somehow failed or key already exists.
	err := rh.APIHandler.Create(c, req.Key, req.Value)
	if err != nil {
		// If key already exists, respond with status code 409.
		if errors.Is(err, repository.ErrAlreadyExists) {
			c.JSON(repository.ErrAlreadyExists.Code, gin.H{"error": "key already exists"})
		} else {
			// If further processing failed, log the error and respond with status code 500.
			rh.Logger.Warn("Tarantool failed to store data",
				"key", req.Key,
				"body", req.Value,
			)
			c.String(http.StatusInternalServerError, "Internal server error.")
		}
		return
	}

	// Everything good --- respond with meta information.
	c.JSON(http.StatusCreated, domain.AppPostResponse{
		Message: "created",
		Key:     req.Key,
		Size:    uint(len(req.Value)), // number of bytes in value field
	})
}

// PUT kv/{id} body: {"value": {ARBITRARY JSON}}
// On success responds with some metainfo about the updated enrty.
func (rh *RequestHandler) PUTHandlerFunc(c *gin.Context) {
	// Empty key is forbidden.
	// Everything else is ok?
	key := c.Param("id")

	if key == "" {
		c.String(http.StatusBadRequest, "Missing key")
		return
	}

	// Verify that body is a valid JSON.
	// Using raw type here because we need to make sure there are no unneccessary keys in the body.
	var rawReq map[string]json.RawMessage
	if err := c.ShouldBindJSON(&rawReq); err != nil {
		c.String(http.StatusBadRequest, "error: invalid request format")
		return
	}

	// Check if there is only single `value` field in a given JSON.
	if len(rawReq) != 1 || rawReq["value"] == nil {
		c.String(http.StatusBadRequest, "error: invalid request format")
		return
	}

	// At this point `value` field is certain to contain JSON data.
	// Process it further.
	value := rawReq["value"]

	// Either Tarantool somehow failed or key was not found.
	err := rh.APIHandler.Update(c, key, value)
	if err != nil {
		// If key was not found, respond with status code 404.
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(repository.ErrNotFound.Code, gin.H{"error": "key not found"})
		} else {
			// If further processing failed, log the error and respond with status code 500.
			rh.Logger.Warn("Tarantool failed to update data",
				"key", key,
				"body", value,
			)
			c.String(http.StatusInternalServerError, "Internal server error.")
		}
		return
	}

	// Everything good --- respond with meta information.
	c.JSON(http.StatusOK, domain.AppPutResponse{
		Message: "updated",
		Key:     key,
		Size:    uint(len(value)), // number of bytes in value field
	})
}

// DELETE kv/{key}
// On success responds with some metainfo about the deleted enrty.
func (rh *RequestHandler) DeleteHandlerFunc(c *gin.Context) {
	// Empty key is forbidden.
	// Everything else is ok?
	key := c.Param("id")

	if key == "" {
		c.String(http.StatusBadRequest, "Missing key")
		return
	}

	// Either Tarantool somehow failed or key was not found.
	err := rh.APIHandler.Delete(c, key)
	if err != nil {
		// If key was not found, respond with status code 404.
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(repository.ErrNotFound.Code, gin.H{"error": "key not found"})
		} else {
			// If further processing failed, log the error and respond with status code 500.
			rh.Logger.Warn("Tarantool failed to delete data",
				"key", key,
			)
			c.String(http.StatusInternalServerError, "Internal server error.")
		}
		return
	}

	// Everything good --- respond with meta information.
	c.JSON(http.StatusOK, domain.AppDeleteResponse{
		Message: "deleted",
		Key:     key,
	})
}
