package handlers

import (
	"net/http"
	"strconv"
	"strings"

	hashiraft "github.com/hashicorp/raft"
	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	nodeID string
	raft   *hashiraft.Raft
}

func NewNodeHandler(
	nodeID string,
	raft *hashiraft.Raft,
) *NodeHandler {

	return &NodeHandler{
		nodeID: nodeID,
		raft:   raft,
	}
}

func (h *NodeHandler) Status(c *gin.Context) {

	stats := h.raft.Stats()

	term, _ := strconv.Atoi(
		stats["term"],
	)

	c.JSON(
		http.StatusOK,
		gin.H{
			"id": h.nodeID,

			"role": strings.ToLower(
				h.raft.State().String(),
			),

			"term": term,

			"leader": string(
				h.raft.Leader(),
			),
		},
	)
}