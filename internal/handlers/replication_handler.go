package handlers

import(
	"gokv/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReplicationHandler struct{
	service *services.KVService
}

func NewReplicationHandler(
	service *services.KVService,
) *ReplicationHandler {

	return &ReplicationHandler{
		service: service,
	}
}

type ReplicationRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (h *ReplicationHandler) ReplicatePut(c *gin.Context) {

	var req ReplicationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	h.service.Put(req.Key, req.Value)

	c.JSON(http.StatusOK, gin.H{
		"message": "replicated",
	})
}

func (h *ReplicationHandler) ReplicateDelete(c *gin.Context) {

	key := c.Param("key")

	h.service.Delete(key)

	c.JSON(http.StatusOK, gin.H{
		"message": "replicated",
	})
}