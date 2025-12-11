package main

import (
	"github.com/len4ernova/go_final_project/pkg/db"
	"github.com/len4ernova/go_final_project/pkg/server"
	"github.com/len4ernova/go_final_project/pkg/services"
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
	logger, err := configZap.Build()
	if err != nil {
		logger.Sugar().Fatalf("can't run logger: %v", err)
	}
	defer logger.Sync()

	config, err := services.GetConfig()
	if err != nil {
		logger.Sugar().Fatalf("can't get config: %v", err)
	}

	// БД
	dataBase, err := db.Init(config.DbName)
	//DBplaner, err := sql.Open("sqlite", configюdbName)
	if err != nil {
		logger.Sugar().Fatalf("Didn't open DB: %v", err)
		return
	}
	defer dataBase.Close()

	logger.Sugar().Info("db = ", config.DbName)

	server.RunSrv(logger, &config, dataBase)

}
