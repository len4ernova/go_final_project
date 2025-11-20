package services

import (
	"os"
	"strconv"
)

// PortVal - получить значение порта из переменной окружения или присвоить значение по ум.
func PortVal(name string, dfltPort int) (int, error) {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		return dfltPort, nil
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	}
	return p, nil
}
