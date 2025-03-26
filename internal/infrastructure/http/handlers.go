// Handlers for router.

package http

import (
	"net/http"
	"tarantool-app/internal/usecase"

	"github.com/gin-gonic/gin"
)

// Incapsulates available handlers.
// All of them mirror CRUD operations so are grouped together.
type RequestHandler struct {
	APIHandler *usecase.CRUD
}

func NewRequestHandler(crudHandler *usecase.CRUD) *RequestHandler {
	return &RequestHandler{APIHandler: crudHandler}
}

// GET kv/{id}
func (rh *RequestHandler) GETHandlerFunc(c *gin.Context) {
	key, ok := c.GetQuery("id")

	if ok && key != "" {
		value, err := rh.APIHandler.Read(c, key)
		if err != nil {
			c.String(http.StatusOK, "All good from GET: %v", value)
		}
	}

	return
}

// POST /kv body: {key: "key", "value": {ARBITRARY JSON}}
func (rh *RequestHandler) POSTHandlerFunc(c *gin.Context) {
	c.String(http.StatusOK, "All good from POST")
	return
}

// PUT kv/{id} body: {"value": {ARBITRARY JSON}}
func (rh *RequestHandler) PUTHandlerFunc(c *gin.Context) {
	key, keyOk := c.GetQuery("id")
	value := c.GetStringMap("value")

	if keyOk {
		err := rh.APIHandler.Update(c, key, value)
		if err != nil {
			c.String(http.StatusOK, "All good from PUT: %v", value)
		}
	}

	return
}

// DELETE kv/{key}
func (r *RequestHandler) DeleteHandlerFunc(c *gin.Context) {
	key, ok := c.GetQuery("id")

	if ok && key != "" {
		value, err := r.APIHandler.Delete(c, key)
		if err != nil {
			c.String(http.StatusOK, "All good from DELETE /kv: %v", value)
		}
	}

	return
}
