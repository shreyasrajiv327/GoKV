package handlers

import (
	"net/http"

	"gokv/internal/models"
	"gokv/internal/services"
	"gokv/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type KVHandler struct {
	service *services.KVService
	logger  *zap.Logger
}

func NewKVHandler(service *services.KVService, logger *zap.Logger) *KVHandler {
	return &KVHandler{
		service: service,
		logger:  logger,
	}
}

func (h *KVHandler) Put(c *gin.Context) {
	var req models.PutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	if req.Key == "" {
		utils.Error(c, http.StatusBadRequest, "key is required")
		return
	}

	h.service.Put(req.Key, req.Value)
	utils.Success(c, http.StatusOK, gin.H{"message": "ok"})
}

func (h *KVHandler) Get(c *gin.Context) {
	key := c.Param("key")

	val, ok := h.service.Get(key)
	if !ok {
		utils.Error(c, http.StatusNotFound, "key not found")
		return
	}

	utils.Success(c, http.StatusOK, models.GetResponse{
		Key:   key,
		Value: val,
	})
}

func (h *KVHandler) Delete(c *gin.Context) {
	key := c.Param("key")
	h.service.Delete(key)
	utils.Success(c, http.StatusOK, gin.H{"message": "deleted"})
}