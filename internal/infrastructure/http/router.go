// Router setup and routes definitions.

package http

import (
	"tarantool-app/internal/constants"
	"tarantool-app/internal/infrastructure/logging"

	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	Engine *gin.Engine
}

func NewGinRouter(env string, logger *logging.Logger, rqHandlers *RequestHandler) *GinRouter {
	if env == constants.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a new engine and attach Recovery and Logger middleware.
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(GinLoggerMiddleware(logger))

	// This is just for silicing a warning.
	// Matters only when a proxy server is involved.
	r.SetTrustedProxies(nil)

	setupRoutes(r, rqHandlers)

	return &GinRouter{Engine: r}
}

func (g *GinRouter) Run(addr string) error {
	return g.Engine.Run(addr)
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
