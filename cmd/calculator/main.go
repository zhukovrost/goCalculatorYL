package main

import (
	"goCalculatorYL/internal/config"
	"goCalculatorYL/internal/router"
	"goCalculatorYL/internal/service"

	"log"
	"net/http"
)

func main() {
	COMPUTING_POWER := 5
	TIME_ADDITION_MS := 1000
	TIME_SUBTRACTION_MS := 1000
	TIME_MULTIPLICATIONS_MS := 1000
	TIME_DIVISIONS_MS := 1000

	// Загрузка конфигурации
	cfg, err := config.LoadConfig(8080, COMPUTING_POWER, TIME_ADDITION_MS, TIME_SUBTRACTION_MS, TIME_MULTIPLICATIONS_MS, TIME_DIVISIONS_MS)
	if err != nil {
		log.Fatalf("Could not load config: %s\n", err.Error())
		return
	}
	// новый service
	srv := service.New(cfg)

	// Настройка маршрутизатора
	r := router.SetupRouter(srv)

	srv.Logger.Infof("Starting server on %s...", cfg.GetAddress())
	// Запуск сервера
	if err := http.ListenAndServe(cfg.GetAddress(), r); err != nil {
		srv.Logger.Fatalf("Could not start server: %s\n", err.Error())
		return
	}
}
