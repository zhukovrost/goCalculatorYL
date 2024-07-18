package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Host           string `json:"address"`
	Port           uint16 `json:"port"`
	Addition       uint
	Subtraction    uint
	Multiplication uint
	Division       uint
}

// GetAddress возвращает полный адрес
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LoadConfig принимает порт для сервера и длительность математических операций и возвращает конфиг
func LoadConfig() *Config {
	addition, subtraction, multiplication, division := loadEnv()
	return &Config{
		Host:           "localhost",
		Port:           8080,
		Addition:       addition,
		Subtraction:    subtraction,
		Multiplication: multiplication,
		Division:       division,
	}
}

func loadEnv() (uint, uint, uint, uint) {
	// Загрузка значений из переменных окружения
	addition := getEnvUint("MATH_ADDITION", 1000)
	subtraction := getEnvUint("MATH_SUBTRACTION", 1000)
	multiplication := getEnvUint("MATH_MULTIPLICATION", 1000)
	division := getEnvUint("MATH_DIVISION", 1000)

	return addition, subtraction, multiplication, division
}

func getEnvUint(key string, defaultValue uint) uint {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid value for %s: %s", key, value)
	}

	if i < 0 {
		log.Fatalf("Invalid value for %s: %s", key, value)
	}

	return uint(i)
}
