package config

import (
	"errors"
	"fmt"
	"time"
)

type Config struct {
	Address                 string        `json:"address"`
	Port                    uint16        `json:"port"`
	TIME_ADDITION_MS        time.Duration `json:"time_addition_ms"`
	TIME_SUBTRACTION_MS     time.Duration `json:"time_subtraction_ms"`
	TIME_MULTIPLICATIONS_MS time.Duration `json:"time_multiplications_ms"`
	TIME_DIVISIONS_MS       time.Duration `json:"time_divisions_ms"`
	COMPUTING_POWER         uint16
}

// GetAddress возвращает полный адрес
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

// LoadConfig принимает порт для сервера и длительность математических операций и возвращает конфиг
func LoadConfig(port uint16, addition, subtraction, multiplication, division, power int) (*Config, error) {
	if port <= 0 || port > 65535 {
		return nil, errors.New("invalid port")
	}

	if addition <= 0 || subtraction <= 0 || multiplication <= 0 || division <= 0 || power <= 0 {
		return nil, errors.New("invalid duration")
	}

	return &Config{
		Address:                 "http://localhost",
		Port:                    port,
		TIME_ADDITION_MS:        time.Millisecond * time.Duration(addition),
		TIME_SUBTRACTION_MS:     time.Millisecond * time.Duration(subtraction),
		TIME_MULTIPLICATIONS_MS: time.Millisecond * time.Duration(multiplication),
		TIME_DIVISIONS_MS:       time.Millisecond * time.Duration(division),
		COMPUTING_POWER:         uint16(power),
	}, nil
}
