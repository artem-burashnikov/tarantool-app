package v1

import (
	"tarantool-app/internal/interfaces"

	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	Engine *gin.Engine
}

func NewGinRouter(env string, log interfaces.Logger, h interfaces.KVHandler) *GinRouter {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Debug("Error setting SetTrustedProxies to nil",
			"error", err,
		)
	}

	setupRoutes(r, h)

	return &GinRouter{Engine: r}
}

func (g *GinRouter) Run(addr string) error {
	return g.Engine.Run(addr)
}

func setupRoutes(r *gin.Engine, h interfaces.KVHandler) {
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	appGroup := r.Group("/kv")
	{
		appGroup.POST("", h.PostKV)
		appGroup.PUT("/:id", h.PutKV)
		appGroup.GET("/:id", h.GetKV)
		appGroup.DELETE("/:id", h.DeleteKV)
	}
}
