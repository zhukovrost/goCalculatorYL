package main

import (
	"goCalculatorYL/internal/app"
	"goCalculatorYL/internal/config"
	"goCalculatorYL/internal/service"
)

func main() {
	const (
		TIME_ADDITION_MS        = 1000
		TIME_SUBTRACTION_MS     = 1000
		TIME_MULTIPLICATIONS_MS = 1000
		TIME_DIVISIONS_MS       = 1000
		PORT                    = 8080
	)
	logger := config.LoadLogger()
	// Загрузка конфигурации
	cfg, err := config.LoadConfig(PORT, TIME_ADDITION_MS, TIME_SUBTRACTION_MS, TIME_MULTIPLICATIONS_MS, TIME_DIVISIONS_MS)
	if err != nil {
		logger.Fatalf("Could not load config: %s\n", err.Error())
		return
	}
	// новый service
	srv := service.NewService(cfg, logger)
	app.Run(srv)
}
