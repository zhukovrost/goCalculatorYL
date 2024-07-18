package main

import (
	"orchestrator/internal/app"
	"orchestrator/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	app.Run(cfg) // запуск приложения
}
