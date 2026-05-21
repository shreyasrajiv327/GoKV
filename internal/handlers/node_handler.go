package handlers

import (
	"net/http"

	"gokv/internal/cluster"

	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	node *cluster.Node
}

func NewNodeHandler(
	node *cluster.Node,
) *NodeHandler {

	return &NodeHandler{
		node: node,
	}
}

func (h *NodeHandler) Status(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"id": h.node.ID,

		"role": h.node.GetRole(),

		"term": h.node.GetTerm(),

		"heartbeat_seconds":
			h.node.TimeSinceHeartbeat().Seconds(),
	})
}