package main

import (
	raftpkg "gokv/internal/raft"

	"gokv/internal/cache"
	"gokv/internal/config"
	"gokv/internal/handlers"
	"gokv/internal/logger"
	"gokv/internal/middleware"
	"gokv/internal/repository"
	"gokv/internal/routes"
	"gokv/internal/services"
	"gokv/internal/snapshot"
	"gokv/internal/wal"

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

	// Memory store
	mem := cache.NewMemoryStore()

	// ---------------- RAFT SETUP ----------------

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
		log.Fatal("failed to setup raft", zap.Error(err))
	}

	_ = raftNode

	// --------------------------------------------

	// Repository
	repo := repository.NewKVRepositroy(mem)

	// Snapshot store
	snapshotStore := snapshot.New("snapshot.json")

	// Load snapshot
	snapshotData, err := snapshotStore.Load()
	if err != nil {
		log.Fatal("failed to load snapshot", zap.Error(err))
	}

	for k, v := range snapshotData {
		mem.Put(k, v)
	}

	// WAL
	walStore, err := wal.New("gokv.wal")
	if err != nil {
		log.Fatal("failed to initialize WAL", zap.Error(err))
	}

	// Replay WAL
	err = walStore.Replay(
		repo.Put,
		repo.Delete,
	)

	if err != nil {
		log.Fatal("failed to replay WAL", zap.Error(err))
	}

	// Background snapshot worker
	go func() {

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {

			err := snapshotStore.Save(mem.GetAll())
			if err != nil {
				log.Error("failed to save snapshot", zap.Error(err))
				continue
			}

			err = walStore.Clear()
			if err != nil {
				log.Error("failed to clear WAL", zap.Error(err))
				continue
			}

			log.Info("snapshot created and WAL compacted")
		}
	}()

	// Service layer
	svc := services.NewKVService(
		repo,
		walStore,
		nil,
		false,
	)

	// Handlers
	h := handlers.NewKVHandler(svc, log)

	// Gin router
	r := gin.New()

	r.Use(middleware.RequestLogger(log))
	r.Use(middleware.Recovery(log))

	// Routes
	routes.Register(r, h)

	log.Info(
		"starting GoKV server",
		zap.String("port", cfg.Port),
		zap.String("node_id", nodeID),
		zap.String("raft_port", raftPort),
	)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("server failed", zap.Error(err))
	}
}