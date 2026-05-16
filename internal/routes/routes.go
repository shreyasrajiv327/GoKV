package routes

import (
	"gokv/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Register(
	r *gin.Engine,
	h *handlers.KVHandler,
) {

	r.PUT("/kv", h.Put)
	r.GET("/kv/:key", h.Get)
	r.DELETE("/kv/:key", h.Delete)
}