// Handlers for router.

package http

import (
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
// Check if id is ok. If it is --- fetch the data. Then return it to user.
func (rh *RequestHandler) GETHandlerFunc(c *gin.Context) {
	key := c.Param("id")

	// Empty key is forbidden.
	if key == "" {
		c.String(http.StatusBadRequest, "Missing key")
		return
	}

	// Any other key is ok. Go fetch the data from repository.
	value, err := rh.APIHandler.Read(c, key)
	// Either Tarantool somehow failed or key was not found.
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(repository.ErrNotFound.Code, gin.H{"error": "key not found"})
		} else {
			rh.Logger.Warn("Tarantool failed to retreive data by key",
				"key", key,
			)
			c.String(http.StatusInternalServerError, "Internal server error.")
		}
		return
	}

	// Key was found --- return in JSON.
	c.JSON(http.StatusOK, domain.AppGetResponse{
		Value: value},
	)
}

// POST /kv body: {key: "key", "value": {ARBITRARY JSON}}
// Verify that body is a valid JSON.
// If so, let Tarantool create an entry.
// Returns some metainfo about the inserted enrty.
func (rh *RequestHandler) POSTHandlerFunc(c *gin.Context) {
	var req domain.AppPostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "error: invalid request format")
		return
	}

	// Either Tarantool somehow failed or key already exists.
	err := rh.APIHandler.Create(c, req.Key, req.Value)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			c.JSON(repository.ErrAlreadyExists.Code, gin.H{"error": "key already exists"})
		} else {
			rh.Logger.Warn("Tarantool failed to store data",
				"key", req.Key,
				"body", req.Value,
			)
			c.String(http.StatusInternalServerError, "Internal server error.")
		}
		return
	}

	// Return meta information to user.
	c.JSON(http.StatusCreated, domain.AppPostResponse{
		Message: "created",
		Key:     req.Key,
		Size:    uint(len(req.Value)), // number of bytes in value field
	})
}

// PUT kv/{id} body: {"value": {ARBITRARY JSON}}
func (rh *RequestHandler) PUTHandlerFunc(c *gin.Context) {
	c.String(http.StatusOK, "All good from PUT")
	return
}

// DELETE kv/{key}
func (r *RequestHandler) DeleteHandlerFunc(c *gin.Context) {
	c.String(http.StatusOK, "All good from DELETE")
	return
}
