package main

import (
	"database/sql"

	"github.com/len4ernova/go_final_project/pkg/db"
	"github.com/len4ernova/go_final_project/pkg/server"
	"github.com/len4ernova/go_final_project/pkg/services"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const dfltPort = 7540
const dfltIp = "127.0.0.1"
const dfltDB = "scheduler.db"

var PlanDB *sql.DB

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

	dbName, err := services.Path2DB("TODO_DBFILE")
	if err != nil {
		logger.Sugar().Fatalf("wrong TODO_DBFILE value")
		return
	}
	if dbName == "" {
		dbName = dfltDB
	}

	err = db.Init(dbName, PlanDB)
	if err != nil {
		logger.Sugar().Fatalf("DB initialization error: %v", err)
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
