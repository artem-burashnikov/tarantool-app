// Handlers for router.

package http

import (
	"errors"
	"net/http"
	"tarantool-app/internal/domain"
	"tarantool-app/internal/interfaces"
	"tarantool-app/internal/repository"

	"github.com/gin-gonic/gin"
)

type AppHandler struct {
	Handler interfaces.UserUseCase
	Logger  interfaces.Logger
}

var _ interfaces.KVHandler = (*AppHandler)(nil) // AppHandler must satisfy KVHandler

func NewRequestHandler(uc interfaces.UserUseCase, log interfaces.Logger) AppHandler {
	return AppHandler{Handler: uc, Logger: log}
}

// GET kv/{id}
func (rh *AppHandler) GetKV(c *gin.Context) {
	rq := domain.Payload{Key: c.Param("id")}

	resp, err := rh.Handler.Read(rq)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrNotFound.Error()})
		} else {
			rh.Logger.Warn("Tarantool failed to retreive data by key",
				"key", rq.Key,
				err,
			)
			c.String(http.StatusInternalServerError, "500 Internal server error.")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":   resp.Key,
		"value": resp.Value,
	})
	return
}

// POST /kv body: {key: "key", "value": {ARBITRARY JSON}}.
func (rh *AppHandler) PostKV(c *gin.Context) {
	var rq domain.Payload

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

	err := rh.Handler.Create(rq)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": repository.ErrAlreadyExists.Error()})
		} else {
			rh.Logger.Warn("Tarantool failed to store data",
				"key", rq.Key,
				"value", rq.Value,
				err,
			)
			c.String(http.StatusInternalServerError, "500 Internal server error.")
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "created",
		"key":     rq.Key,
		"value":   rq.Value,
	})
	return
}

// PUT kv/{id} body: {"value": {ARBITRARY JSON}}
func (rh *AppHandler) PutKV(c *gin.Context) {
	rq := domain.Payload{Key: c.Param("id")}

	if err := c.ShouldBindJSON(&rq.Value); err != nil || len(rq.Value) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	if rq.Value["value"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing value"})
		return
	}

	err := rh.Handler.Update(rq)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrNotFound.Error()})
		} else {
			rh.Logger.Warn("Tarantool failed to update data",
				"key", rq.Key,
				"body", rq.Value,
				err,
			)
			c.String(http.StatusInternalServerError, "500 Internal server error.")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "updated",
		"key":     rq.Key,
		"value":   rq.Value,
	})
	return
}

// DELETE kv/{key}
func (rh *AppHandler) DeleteKV(c *gin.Context) {
	rq := domain.Payload{Key: c.Param("id")}

	resp, err := rh.Handler.Delete(rq)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrNotFound.Error()})
		} else {
			rh.Logger.Warn("Tarantool failed to delete data",
				"key", rq.Key,
				err,
			)
			c.String(http.StatusInternalServerError, "500 Internal server error.")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
		"key":     resp.Key,
		"value":   resp.Value,
	})
	return
}
