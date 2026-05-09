package main

import (
	"gokv/internal/cache"
	"gokv/internal/config"
	"gokv/internal/handlers"
	"gokv/internal/logger"
	"gokv/internal/repository"
	"gokv/internal/routes"
	"gokv/internal/services"
    "gokv/internal/middleware"
	"gokv/internal/wal"
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

	mem := cache.NewMemoryStore()
	repo := repository.NewKVRepositroy(mem)
	walStore, err := wal.New("gokv.wal")
	if err != nil{
		log.Fatal("failed to initialize WAL", zap.Error(err))
	}

	err = walStore.Replay(
		repo.Put,
		repo.Delete,
	)

	if err != nil{
		log.Fatal("failed to replay WAL", zap.Error(err))
	}
	svc := services.NewKVService(repo, walStore)
	h := handlers.NewKVHandler(svc, log)

	r := gin.New()
	r.Use(middleware.RequestLogger(log))
	r.Use(middleware.Recovery(log))

	routes.Register(r, h)

	log.Info("starting GoKV server", zap.String("port", cfg.Port))

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("server failed", zap.Error(err))
	}
}