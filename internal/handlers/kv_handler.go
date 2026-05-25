package handlers

import (
	"encoding/json"
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

func NewKVHandler(
	service *services.KVService,
	logger *zap.Logger,
) *KVHandler {

	return &KVHandler{
		service: service,
		logger:  logger,
	}
}

func (h *KVHandler) Put(
	c *gin.Context,
) {

	var req models.PutRequest

	if err :=
		c.ShouldBindJSON(
			&req,
		); err != nil {

		utils.Error(
			c,
			http.StatusBadRequest,
			"invalid request",
		)

		return
	}

	if req.Key == "" {

		utils.Error(
			c,
			http.StatusBadRequest,
			"key required",
		)

		return
	}

	err :=
		h.service.Put(
			req.Key,
			req.Value,
		)

	if err != nil {

		utils.Error(
			c,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	utils.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "ok",
		},
	)
}

func (h *KVHandler) Get(
	c *gin.Context,
) {

	key :=
		c.Param("key")

	consistency :=
		c.DefaultQuery(
			"consistency",
			"eventual",
		)

	if consistency == "smart" {

		if !h.service.CanServeSmartRead() {

			resp,
				err :=
				h.service.
					ForwardStrongRead(
						key,
					)

			if err != nil {

				utils.Error(
					c,
					http.StatusServiceUnavailable,
					err.Error(),
				)

				return
			}

			defer resp.Body.Close()

			var body map[string]interface{}

			json.NewDecoder(
				resp.Body,
			).Decode(
				&body,
			)

			c.JSON(
				resp.StatusCode,
				body,
			)

			return
		}
	}

	val,
		ok,
		err :=
		h.service.Get(
			key,
			consistency,
		)

	if err != nil {

		utils.Error(
			c,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	if !ok {

		utils.Error(
			c,
			http.StatusNotFound,
			"key not found",
		)

		return
	}

	utils.Success(
		c,
		http.StatusOK,
		models.GetResponse{
			Key:   key,
			Value: val,
		},
	)
}

func (h *KVHandler) Delete(
	c *gin.Context,
) {

	key :=
		c.Param("key")

	err :=
		h.service.Delete(
			key,
		)

	if err != nil {

		utils.Error(
			c,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	utils.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "deleted",
		},
	)
}