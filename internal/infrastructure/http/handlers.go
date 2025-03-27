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
	rq := domain.AppGetRequest{Key: c.Param("id")}

	// Empty key is forbidden.
	// Everything else is ok?
	if rq.Key == "" {
		c.String(http.StatusBadRequest, "Missing key")
		return
	}

	// Any other key is ok. Go fetch the data from repository.
	resp, err := rh.APIHandler.Read(c, &rq)
	// Either Tarantool somehow failed or key was not found.
	// TODO: " " key returns 500 internal server errror
	if err != nil {
		// If key was not found, respond with status code 404.
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrNotFound.Message})
		} else {
			// If further processing failed, log the error and respond with status code 500.
			rh.Logger.Warn("Tarantool failed to retreive data by key",
				"key", rq.Key,
			)
			c.String(http.StatusInternalServerError, "Internal server error.")
		}
		return
	}

	// Key was found --- respond in JSON.
	c.JSON(http.StatusOK, resp)
	return
}

// POST /kv body: {key: "key", "value": {ARBITRARY JSON}}.
// On success responds with some metainfo about the inserted enty.
func (rh *RequestHandler) POSTHandlerFunc(c *gin.Context) {
	// Verify that body is a valid JSON.
	// If so, process it further and let Tarantool create an entry.
	// Incorrect request body returns status code 400.
	var req domain.AppPostRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.Key == "" || len(req.Value) == 0 {
		c.String(http.StatusBadRequest, "error: invalid request format")
		return
	}

	// Either Tarantool somehow failed or key already exists.
	resp, err := rh.APIHandler.Create(c, &req)
	if err != nil {
		// If key already exists, respond with status code 409.
		if errors.Is(err, repository.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": repository.ErrAlreadyExists.Message})
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

	// Everything good --- respond.
	c.JSON(http.StatusCreated, resp)
	return
}

// PUT kv/{id} body: {"value": {ARBITRARY JSON}}
// On success responds with some metainfo about the updated enty.
func (rh *RequestHandler) PUTHandlerFunc(c *gin.Context) {
	rq := domain.AppPutRequest{Key: c.Param("id")}

	// Empty key is forbidden.
	// Everything else is ok?
	if rq.Key == "" {
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
	rq.Value = rawReq["value"]

	// Either Tarantool somehow failed or key was not found.
	resp, err := rh.APIHandler.Update(c, &rq)
	if err != nil {
		// If key was not found, respond with status code 404.
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrNotFound.Message})
		} else {
			// If further processing failed, log the error and respond with status code 500.
			rh.Logger.Warn("Tarantool failed to update data",
				"key", rq.Key,
				"body", rq.Value,
			)
			c.String(http.StatusInternalServerError, "Internal server error.")
		}
		return
	}

	// Everything good --- respond.
	c.JSON(http.StatusOK, resp)
	return
}

// DELETE kv/{key}
// On success responds with some metainfo about the deleted enty.
func (rh *RequestHandler) DeleteHandlerFunc(c *gin.Context) {
	rq := domain.AppDeleteRequest{Key: c.Param("id")}

	// Empty key is forbidden.
	// Everything else is ok?
	if rq.Key == "" {
		c.String(http.StatusBadRequest, "Missing key")
		return
	}

	// Either Tarantool somehow failed or key was not found.
	resp, err := rh.APIHandler.Delete(c, &rq)
	if err != nil {
		// If key was not found, respond with status code 404.
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrNotFound.Message})
		} else {
			// If further processing failed, log the error and respond with status code 500.
			rh.Logger.Warn("Tarantool failed to delete data",
				"key", rq.Key,
			)
			c.String(http.StatusInternalServerError, "Internal server error.")
		}
		return
	}

	// Everything good --- respond.
	c.JSON(http.StatusOK, resp)
	return
}
