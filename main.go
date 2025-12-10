package main

import (
	"database/sql"
	"os"

	"github.com/len4ernova/go_final_project/pkg/db"
	"github.com/len4ernova/go_final_project/pkg/server"
	"github.com/len4ernova/go_final_project/pkg/services"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const dfltPort = 7540
const dfltDB = "scheduler.db"

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
		logger.Sugar().Fatalf("wrong TODO_PORT value")
		return
	}

	dbName := os.Getenv("TODO_DBFILE")
	if dbName == "" {
		dbName = dfltDB
	}

	// БД
	DBplaner, err := sql.Open("sqlite", dbName)
	if err != nil {
		logger.Fatal("Didn't open DB.")
		return
	}
	defer DBplaner.Close()
	err = db.Init(dbName, DBplaner)
	if err != nil {
		logger.Sugar().Fatalf("DB initialization error: %v", err)
	}

	// настройки сервера
	settingsSrv := server.Settings{
		// Ip:   dfltIp,
		Port: port,
	}

	logger.Sugar().Info("db = ", dbName)
	server.RunSrv(logger, &settingsSrv, DBplaner)
	// if err != nil {
	// 	logger.Fatal(err.Error())
	// }
}
