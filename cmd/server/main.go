package main

import (
	"os"

	raftpkg "gokv/internal/raft"

	"gokv/internal/cache"
	"gokv/internal/config"
	"gokv/internal/handlers"
	"gokv/internal/logger"
	"gokv/internal/middleware"
	"gokv/internal/repository"
	"gokv/internal/routes"
	"gokv/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	cfg := config.Load()

	log, err := logger.New()
	if err != nil {
		panic(err)
	}

	defer log.Sync()

	//--------------------------------
	// Memory Store
	//--------------------------------

	mem := cache.NewMemoryStore()
	repo :=
		repository.
			NewKVRepositroy(
				mem,
			)

	//--------------------------------
	// RAFT
	//--------------------------------

	nodeID := os.Getenv("NODE_ID")
	raftPort := os.Getenv("RAFT_PORT")

	raftDir :=
		"raft-data-" +
			nodeID

	raftNode,
		fsm,
		err :=
		raftpkg.SetupRaft(
			nodeID,
			raftDir,
			"127.0.0.1:"+raftPort,
			repo,
		)

	if err != nil {

		log.Fatal(
			"failed to setup raft",
			zap.Error(err),
		)
	}

	//--------------------------------
	// HTTP Node Mapping
	//--------------------------------

	httpNodes :=
		map[string]string{

			"127.0.0.1:7001": "8001",

			"127.0.0.1:7002": "8002",

			"127.0.0.1:7003": "8003",
		}

	//--------------------------------
	// Service Layer
	//--------------------------------

	svc :=
		services.
			NewKVService(
				repo,
				raftNode,
				fsm,
				httpNodes,
			)

	//--------------------------------
	// Handlers
	//--------------------------------

	kvHandler :=
		handlers.
			NewKVHandler(
				svc,
				log,
			)

	nodeHandler :=
		handlers.
			NewNodeHandler(
				nodeID,
				raftNode,
			)

	//--------------------------------
	// Router
	//--------------------------------

	r := gin.New()

	r.Use(
		middleware.
			RequestLogger(log),
	)

	r.Use(
		middleware.
			Recovery(log),
	)

	routes.Register(
		r,
		kvHandler,
		nodeHandler,
	)

	log.Info(
		"starting GoKV",

		zap.String(
			"port",
			cfg.Port,
		),

		zap.String(
			"node",
			nodeID,
		),

		zap.String(
			"raft",
			raftPort,
		),
	)

	err =
		r.Run(
			":" +
				cfg.Port,
		)

	if err != nil {

		log.Fatal(
			"server failed",
			zap.Error(err),
		)
	}
}
