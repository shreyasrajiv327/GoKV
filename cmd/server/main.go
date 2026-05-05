// cmd/server/main.go
package main

import (
	"gokv/internal/cache"
	"gokv/internal/config"
	"gokv/internal/handlers"
	"gokv/internal/logger"
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

	mem := cache.NewMemoryStore()
	repo := repository.NewKVRepositroy(mem)
	svc := services.NewKVService(repo)
	h := handlers.NewKVHandler(svc, log)

	r := gin.New()
	r.Use(gin.Recovery())

	routes.Register(r, h)

	log.Info("starting GoKV server", zap.String("port", cfg.Port))

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("server failed", zap.Error(err))
	}
}