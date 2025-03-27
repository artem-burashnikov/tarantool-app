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
// On success responds with value stored in database.
func (rh *RequestHandler) GETHandlerFunc(c *gin.Context) {
	rq := domain.AppPack{Key: c.Param("id")}

	resp, err := rh.APIHandler.Read(c, &rq)
	// Either Tarantool somehow failed or key was not found.
	if err != nil {
		// If key was not found, respond with status code 404.
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrNotFound.Error()})
		} else {
			// If further processing failed, log the error and respond with status code 500.
			rh.Logger.Warn("Tarantool failed to retreive data by key",
				"key", rq.Key,
				err,
			)
			c.String(http.StatusInternalServerError, "500 Internal server error.")
		}
		return
	}

	// Key was found and request now contains repository response --- respond to user.
	c.JSON(http.StatusOK, gin.H{
		"key":   resp.Key,
		"value": resp.Value,
	})
	return
}

// POST /kv body: {key: "key", "value": {ARBITRARY JSON}}.
// On success responds with the inserted entry.
func (rh *RequestHandler) POSTHandlerFunc(c *gin.Context) {
	var rq domain.AppPack

	// Validate the request body.
	// If ok, let Tarantool create an entry.
	// Incorrect request body returns status code 400.
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	if rq.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing key"})
		return
	}

	if len(rq.Value) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing value"})
		return
	}

	// Either Tarantool somehow failed or key already exists.
	err := rh.APIHandler.Create(c, &rq)
	if err != nil {
		// If key already exists, respond with status code 409.
		if errors.Is(err, repository.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": repository.ErrAlreadyExists.Error()})
		} else {
			// If further processing failed, log the error and respond with status code 500.
			rh.Logger.Warn("Tarantool failed to store data",
				"key", rq.Key,
				"value", rq.Value,
				err,
			)
			c.String(http.StatusInternalServerError, "500 Internal server error.")
		}
		return
	}

	// Everything good --- respond.
	c.JSON(http.StatusCreated, gin.H{
		"message": "created",
		"key":     rq.Key,
		"value":   rq.Value,
	})
	return
}

// PUT kv/{id} body: {"value": {ARBITRARY JSON}}
// On success responds with the updated enty.
func (rh *RequestHandler) PUTHandlerFunc(c *gin.Context) {
	rq := domain.AppPack{Key: c.Param("id")}

	// Validate the body.
	if err := c.ShouldBindJSON(&rq.Value); err != nil || len(rq.Value) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	if rq.Value["value"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing value"})
		return
	}

	// Process the request.
	err := rh.APIHandler.Update(c, &rq)
	// Either Tarantool somehow failed or key was not found.
	if err != nil {
		// If key was not found, respond with status code 404.
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrNotFound.Error()})
		} else {
			// If further processing failed, log the error and respond with status code 500.
			rh.Logger.Warn("Tarantool failed to update data",
				"key", rq.Key,
				"body", rq.Value,
				err,
			)
			c.String(http.StatusInternalServerError, "500 Internal server error.")
		}
		return
	}

	// Everything good --- respond.
	c.JSON(http.StatusOK, gin.H{
		"message": "updated",
		"key":     rq.Key,
		"value":   rq.Value,
	})
	return
}

// DELETE kv/{key}
// On success responds with the deleted entry.
func (rh *RequestHandler) DeleteHandlerFunc(c *gin.Context) {
	rq := domain.AppPack{Key: c.Param("id")}

	// Either Tarantool somehow failed or key was not found.
	resp, err := rh.APIHandler.Delete(c, &rq)
	if err != nil {
		// If key was not found, respond with status code 404.
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrNotFound.Error()})
		} else {
			// If further processing failed, log the error and respond with status code 500.
			rh.Logger.Warn("Tarantool failed to delete data",
				"key", rq.Key,
				err,
			)
			c.String(http.StatusInternalServerError, "500 Internal server error.")
		}
		return
	}

	// Everything good --- respond.
	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
		"key":     resp.Key,
		"value":   resp.Value,
	})
	return
}
