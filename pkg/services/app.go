package services

import (
	"os"
	"strconv"
)

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

// Path2DB - получение значения путь к БД.
func Path2DB(envName string) (string, error) {
	val := os.Getenv(envName)
	//TODO validator
	return val, nil
}
