package main

import (
	"goCalculatorYL/internal/config"
	"goCalculatorYL/internal/router"
	"goCalculatorYL/internal/service"
	"log"
	"net/http"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig(8080, 100, 100, 100, 100, 100)
	if err != nil {
		log.Fatalf("Could not load config: %s\n", err.Error())
	}
	srv := service.New(cfg)

	// Настройка маршрутизатора
	r := router.SetupRouter(srv)

	// Запуск сервера
	log.Printf("Starting server on %s...", cfg.GetAddress())
	if err := http.ListenAndServe(cfg.GetAddress(), r); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
