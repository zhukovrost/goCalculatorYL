package app

import (
	"net/http"
	"orchestrator/internal/config"
	"orchestrator/internal/handlers"
	"orchestrator/internal/router"
	"orchestrator/internal/service"
	"orchestrator/pkg/logger"
	"orchestrator/pkg/sqlite"
)

func Run(cfg *config.Config) {
	log := logger.New(true)
	log.Info("Starting orchestrator...")

	db, err := sqlite.Open()
	if err != nil {
		log.Fatalf("Failed to connect to database: %s\n", err.Error())
		return
	}

	defer db.Close()

	srv := service.New(cfg, db, log)
	handler := handlers.New(srv)

	// Настройка маршрутизатора
	r := router.SetupRouter(handler)

	srv.Logger.Infof("Starting server on %s...", srv.Cfg.GetAddress())
	// Запуск сервера
	if err := http.ListenAndServe(srv.Cfg.GetAddress(), r); err != nil {
		srv.Logger.Fatalf("Could not start server: %s\n", err.Error())
		return
	}
}
