package utils

import (
	"gokv/internal/models"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, models.ErrorResponse{
		Error: message,
	})
}