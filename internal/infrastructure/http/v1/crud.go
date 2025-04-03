// Handlers for router.

package v1

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

var _ interfaces.KVHandler = AppHandler{} // AppHandler must satisfy KVHandler

func NewRequestHandler(uc interfaces.UserUseCase, log interfaces.Logger) AppHandler {
	return AppHandler{Handler: uc, Logger: log}
}

// @Summary      Get value by key
// @Description  Retrieves the value for the specified key from the Tarantool database.
// @Tags         kv
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Key ID"
// @Success      200 {object} map[string]interface{} "Success"
// @Failure      404 {object} map[string]interface{} "Key not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /kv/{id} [get]
func (rh AppHandler) GetKV(c *gin.Context) {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "500 Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":   resp.Key,
		"value": resp.Value,
	})
	return //nolint:staticcheck
}

// @Summary      Create a new key-value pair
// @Description  Creates a new key with the provided value in the Tarantool database.
// @Tags         kv
// @Accept       json
// @Produce      json
// @Param        body  body  domain.Payload  true  "Payload containing key and value"
// @Success      201 {object} map[string]interface{} "Created successfully"
// @Failure      400 {object} map[string]interface{} "Invalid request"
// @Failure      409 {object} map[string]interface{} "Key already exists"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /kv [post]
func (rh AppHandler) PostKV(c *gin.Context) {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "500 Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "created",
		"key":     rq.Key,
		"value":   rq.Value,
	})
	return //nolint:staticcheck
}

// @Summary      Update value by key
// @Description  Updates the value for the specified key in the Tarantool database.
// @Tags         kv
// @Accept       json
// @Produce      json
// @Param        id    path  string          true  "Key ID"
// @Param        body  body  domain.Payload  true  "Payload containing updated value"
// @Success      200 {object} map[string]interface{} "Updated successfully"
// @Failure      400 {object} map[string]interface{} "Invalid request"
// @Failure      404 {object} map[string]interface{} "Key not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /kv/{id} [put]
func (rh AppHandler) PutKV(c *gin.Context) {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "500 Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "updated",
		"key":     rq.Key,
		"value":   rq.Value,
	})
	return //nolint:staticcheck
}

// @Summary      Delete key-value pair
// @Description  Deletes the specified key and its value from the Tarantool database.
// @Tags         kv
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Key ID"
// @Success      200 {object} map[string]interface{} "Deleted successfully"
// @Failure      404 {object} map[string]interface{} "Key not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /kv/{id} [delete]
func (rh AppHandler) DeleteKV(c *gin.Context) {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "500 Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
		"key":     resp.Key,
		"value":   resp.Value,
	})
	return //nolint:staticcheck
}
