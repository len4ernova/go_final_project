package server

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/len4ernova/go_final_project/pkg/api"
	"go.uber.org/zap"
)

type Settings struct {
	Ip   string
	Port int
}

// RunSrv - запустить сервер.
func RunSrv(logger *zap.Logger, settings *Settings, dbPlan *sql.DB) {
	// структура с логгером, обработчиками
	hands := api.SrvHand{
		Logger: logger,
		DB:     dbPlan,
	}

	// создание своего сервера
	mux := http.NewServeMux()
	server := http.Server{
		Addr:        settings.Ip + ":" + strconv.Itoa(settings.Port),
		Handler:     mux,
		ErrorLog:    nil,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 10 * time.Second,
	}

	// endpoints
	mux.Handle("GET /css/", http.StripPrefix("/css", http.FileServer(http.Dir("./web/css"))))
	mux.Handle("GET /js/", http.StripPrefix("/js", http.FileServer(http.Dir("./web/js"))))
	mux.Handle("/", http.FileServer(http.Dir("./web/")))

	mux.HandleFunc("GET /{$}", hands.Index)

	hands.Init(mux)

	logger.Sugar().Info("Serving on http://%v ...", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Sugar().Fatal(err)
	}
}
