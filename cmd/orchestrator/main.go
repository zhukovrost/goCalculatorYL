package main

import (
	"goCalculatorYL/internal/app"
	"goCalculatorYL/internal/config"
	"goCalculatorYL/internal/service"
)

func main() {
	// настройки
	const (
		TIME_ADDITION_MS        = 1500
		TIME_SUBTRACTION_MS     = 2000
		TIME_MULTIPLICATIONS_MS = 5000
		TIME_DIVISIONS_MS       = 4000
		PORT                    = 8080
		DEBUG_LEVEL             = true
	)

	logger := config.LoadLogger(DEBUG_LEVEL) // загрузка логгера
	cfg, err := config.LoadConfig(PORT, TIME_ADDITION_MS, TIME_SUBTRACTION_MS,
		TIME_MULTIPLICATIONS_MS, TIME_DIVISIONS_MS) // Загрузка конфигурации
	if err != nil {
		logger.Fatalf("Could not load config: %s\n", err.Error())
		return
	}
	srv := service.NewService(cfg, logger) // новый service
	app.Run(srv)                           // запуск приложения
}
