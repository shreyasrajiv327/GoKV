package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"gokv/internal/models"

	hashiraft "github.com/hashicorp/raft"
	"github.com/gin-gonic/gin"
)

type NodeHandler struct{

	nodeID string

	raft *hashiraft.Raft
}

func NewNodeHandler(
	nodeID string,
	raft *hashiraft.Raft,
)*NodeHandler{

	return &NodeHandler{
		nodeID:nodeID,
		raft:raft,
	}
}

func(h *NodeHandler) Status(
	c *gin.Context,
){

	stats:=h.raft.Stats()

	term,_:=
		strconv.Atoi(
			stats["term"],
		)

	c.JSON(
		http.StatusOK,
		gin.H{
			"id":h.nodeID,
			"role":strings.ToLower(
				h.raft.State().String(),
			),
			"term":term,
			"leader":string(
				h.raft.Leader(),
			),
		},
	)
}

func(h *NodeHandler) Join(
	c *gin.Context,
){

	if h.raft.State()!=
		hashiraft.Leader{

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error":"not leader",
			},
		)

		return
	}

	var req models.JoinRequest

	if err:=
		c.ShouldBindJSON(
			&req,
		);err!=nil{

		c.JSON(
			400,
			gin.H{
				"error":"bad request",
			},
		)

		return
	}

	future:=
		h.raft.AddVoter(

			hashiraft.ServerID(
				req.NodeID,
			),

			hashiraft.ServerAddress(
				req.RaftAddr,
			),

			0,
			0,
		)

	if err:=
		future.Error();err!=nil{

		c.JSON(
			500,
			gin.H{
				"error":err.Error(),
			},
		)

		return
	}

	c.JSON(
		200,
		gin.H{
			"message":"joined",
		},
	)
}