package main

import (
	"github.com/len4ernova/go_final_project/pkg/server"
	"github.com/len4ernova/go_final_project/pkg/services"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const dfltPort = 7540
const dfltIp = "127.0.0.1"

func main() {
	// Настройка логгера: вывода логов в консоль в формате JSON
	configZap := zap.Config{
		Encoding:      "json",
		Level:         zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:   []string{"stdout"}, // вывод в консоль // TODO >file
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}
	logger, _ := configZap.Build()
	defer logger.Sync()

	port, err := services.PortVal("TODO_PORT", dfltPort)
	if err != nil {
		logger.Sugar().Fatalf("wrong port value")
		return
	}

	settingsSrv := server.Settings{
		Ip:   dfltIp,
		Port: port,
	}

	server.RunSrv(logger, &settingsSrv)
	// if err != nil {
	// 	logger.Fatal(err.Error())
	// }
}
