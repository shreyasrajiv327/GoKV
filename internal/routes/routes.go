package routes

import (
	"gokv/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Register(
	r *gin.Engine,
	kvHandler *handlers.KVHandler,
	nodeHandler *handlers.NodeHandler,
) {

	// KV APIs
	r.PUT("/kv", kvHandler.Put)
	r.GET("/kv/:key", kvHandler.Get)
	r.DELETE("/kv/:key", kvHandler.Delete)

	// Cluster APIs
	r.GET("/node/status", nodeHandler.Status)
	r.POST(
	"/cluster/join",
	nodeHandler.Join,
)
}