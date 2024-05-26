package main

import (
	"flag"
	"goCalculatorYL/internal/app"
	"goCalculatorYL/internal/config"
	"goCalculatorYL/internal/service"
)

func main() {
	const (
		PORT = 8080
	)

	// Определение флагов (подробнее go run cmd/orchestrator/main.go -h)
	TIME_ADDITION_MS := flag.Uint("add", 5000, "time for addition operation in milliseconds")
	TIME_SUBTRACTION_MS := flag.Uint("sub", 6000, "time for subtraction operation in milliseconds")
	TIME_MULTIPLICATIONS_MS := flag.Uint("mul", 10000, "time for multiplication operation in milliseconds")
	TIME_DIVISIONS_MS := flag.Uint("div", 11000, "time for division operation in milliseconds")
	debug := flag.Bool("debug", false, "enable debug level logging")

	// Парсинг флагов
	flag.Parse()

	logger := config.LoadLogger(*debug) // загрузка логгера
	cfg, err := config.LoadConfig(PORT, *TIME_ADDITION_MS, *TIME_SUBTRACTION_MS,
		*TIME_MULTIPLICATIONS_MS, *TIME_DIVISIONS_MS) // Загрузка конфигурации
	if err != nil {
		logger.Fatalf("Could not load config: %s\n", err.Error())
		return
	}
	srv := service.NewService(cfg, logger) // новый service
	app.Run(srv)                           // запуск приложения
}
