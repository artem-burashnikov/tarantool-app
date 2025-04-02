package interfaces

import "github.com/gin-gonic/gin"

type KVHandler interface {
	GetKV(c *gin.Context)    // GET /kv/:id
	PostKV(c *gin.Context)   // POST /kv
	PutKV(c *gin.Context)    // PUT /kv/:id
	DeleteKV(c *gin.Context) // DELETE /kv/:id
}
