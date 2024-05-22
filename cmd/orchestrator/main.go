package main

import (
	"goCalculatorYL/internal/app"
	"goCalculatorYL/internal/config"
	"goCalculatorYL/internal/service"
)

func main() {
	logger := config.LoadLogger()
	// Загрузка конфигурации
	cfg, err := config.LoadConfig(8080)
	if err != nil {
		logger.Fatalf("Could not load config: %s\n", err.Error())
		return
	}
	// новый service
	srv := service.New(cfg, logger)
	app.Run(srv)
}
