package services

import (
	"fmt"
	"os"
	"strconv"
)

const dfltPort = 7540
const dfltDB = "scheduler.db"

type Config struct {
	Port   int
	DbName string
}

// GetConfig - установка настроек.
func GetConfig() (Config, error) {
	var cfg Config
	var err error

	cfg.Port, err = PortVal("TODO_PORT", dfltPort)
	if err != nil {
		return Config{}, fmt.Errorf("wrong TODO_PORT value: %f", err)
	}

	cfg.DbName = os.Getenv("TODO_DBFILE")
	if cfg.DbName == "" {
		cfg.DbName = dfltDB
	}
	return cfg, nil
}

// PortVal - получение значения порт из переменной окружения или присвоить значение по ум.
func PortVal(envName string, dfltPort int) (int, error) {
	port := os.Getenv(envName)
	if port == "" {
		return dfltPort, nil
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	}
	return p, nil
}
