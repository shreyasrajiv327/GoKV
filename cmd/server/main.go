package main

import (
	raftpkg "gokv/internal/raft"

	"gokv/internal/cache"
	"gokv/internal/cluster"
	"gokv/internal/config"
	"gokv/internal/handlers"
	"gokv/internal/logger"
	"gokv/internal/middleware"
	"gokv/internal/repository"
	"gokv/internal/routes"
	"gokv/internal/services"
	"gokv/internal/snapshot"
	"gokv/internal/wal"

	hashiraft "github.com/hashicorp/raft"

	"os"
	"time"

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

	//--------------------------------
	// RAFT SETUP
	//--------------------------------

	nodeID := os.Getenv("NODE_ID")
	raftPort := os.Getenv("RAFT_PORT")

	raftDir := "raft-data-" + nodeID

	raftNode, err := raftpkg.SetupRaft(
		nodeID,
		raftDir,
		"127.0.0.1:"+raftPort,
		mem,
	)

	if err != nil {

		log.Fatal(
			"failed to setup raft",
			zap.Error(err),
		)
	}

	//--------------------------------
	// Node state (visibility/debug)
	//--------------------------------

	node := cluster.NewNode(
		nodeID,
		cluster.Follower,
	)

	if raftNode.State() == hashiraft.Leader {

		node.SetRole(
			cluster.Leader,
		)
	}

	//--------------------------------
	// Repository
	//--------------------------------

	repo :=
		repository.NewKVRepositroy(
			mem,
		)

	//--------------------------------
	// Snapshot
	//--------------------------------

	snapshotStore :=
		snapshot.New(
			"snapshot.json",
		)

	snapshotData, err :=
		snapshotStore.Load()

	if err != nil {

		log.Fatal(
			"failed snapshot load",
			zap.Error(err),
		)
	}

	for k, v := range snapshotData {

		mem.Put(k, v)
	}

	//--------------------------------
	// WAL
	//--------------------------------

	walStore, err :=
		wal.New(
			"gokv.wal",
		)

	if err != nil {

		log.Fatal(
			"failed WAL init",
			zap.Error(err),
		)
	}

	err =
		walStore.Replay(
			repo.Put,
			repo.Delete,
		)

	if err != nil {

		log.Fatal(
			"WAL replay failed",
			zap.Error(err),
		)
	}

	//--------------------------------
	// Snapshot Worker
	//--------------------------------

	go func() {

		ticker :=
			time.NewTicker(
				30 * time.Second,
			)

		defer ticker.Stop()

		for range ticker.C {

			err :=
				snapshotStore.Save(
					mem.GetAll(),
				)

			if err != nil {

				log.Error(
					"snapshot failed",
					zap.Error(err),
				)

				continue
			}

			err =
				walStore.Clear()

			if err != nil {

				log.Error(
					"WAL clear failed",
					zap.Error(err),
				)

				continue
			}

			log.Info(
				"snapshot complete + WAL compacted",
			)
		}

	}()

	//--------------------------------
	// Service Layer
	//--------------------------------

	svc :=
		services.NewKVService(
			repo,
			walStore,
			nil,
			false,
		)

	//--------------------------------
	// Handlers
	//--------------------------------

	kvHandler :=
		handlers.NewKVHandler(
			svc,
			log,
		)

	nodeHandler :=
		handlers.NewNodeHandler(
			node,
		)

	//--------------------------------
	// Router
	//--------------------------------

	r := gin.New()

	r.Use(
		middleware.RequestLogger(log),
	)

	r.Use(
		middleware.Recovery(log),
	)

	routes.Register(
		r,
		kvHandler,
		nodeHandler,
	)

	//--------------------------------
	// Startup
	//--------------------------------

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
			":" + cfg.Port,
		)

	if err != nil {

		log.Fatal(
			"server failed",
			zap.Error(err),
		)
	}
}