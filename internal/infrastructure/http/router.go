// Router interface and implementation.
// The interface may be used for testing.

package http

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	Run(addr string) error
}

type GinRouter struct {
	engine *gin.Engine
}

func NewGinRouter(rqHandlers *RequestHandler) *GinRouter {
	r := gin.Default()
	setupRoutes(r, rqHandlers)
	return &GinRouter{engine: r}
}

// Satisfies Router interface.
func (g *GinRouter) Run(addr string) error {
	return g.engine.Run(addr)
}

// Binds handlers to route handles.
func setupRoutes(r *gin.Engine, h *RequestHandler) {
	appGroup := r.Group("/kv")
	{
		appGroup.POST("", h.POSTHandlerFunc)
		appGroup.PUT("/:id", h.PUTHandlerFunc)
		appGroup.GET("/:id", h.GETHandlerFunc)
		appGroup.DELETE("/:id", h.DeleteHandlerFunc)
	}
}
