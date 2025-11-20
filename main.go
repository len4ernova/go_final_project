package main

import (
	"github.com/len4ernova/go_final_project/pkg/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

	settings := server.Settings{
		Ip:   "127.0.0.1",
		Port: 7540,
	}

	server.RunSrv(logger, &settings)
	// if err != nil {
	// 	logger.Fatal(err.Error())
	// }
}
