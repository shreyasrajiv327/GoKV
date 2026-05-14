package routes

import (
	"gokv/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, h *handlers.KVHandler, replicationHandler *handlers.ReplicationHandler) {
	r.PUT("/kv", h.Put)
	r.GET("/kv/:key", h.Get)
	r.DELETE("/kv/:key", h.Delete)
	r.POST("/replicate/put", replicationHandler.ReplicatePut)
    r.DELETE("/replicate/delete/:key", replicationHandler.ReplicateDelete)
}