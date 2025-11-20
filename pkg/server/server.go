package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/len4ernova/go_final_project/pkg/handlers"
	"go.uber.org/zap"
)

type Settings struct {
	Ip   string
	Port int
}
const indexHTML = "web/index.html"

const webDir = "../go_final_project/web"
// RunSrv - запустить сервер.
func RunSrv(logger *zap.Logger, settings *Settings) {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:        settings.Ip + ":" + strconv.Itoa(settings.Port),
		Handler:     mux,
		ErrorLog:    nil,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 10 * time.Second,
	}
    
    mux.Handle("GET /css/", http.StripPrefix("/css", http.FileServer(http.Dir("./web/css"))))
    mux.Handle("GET /js/", http.StripPrefix("/js", http.FileServer(http.Dir("./web/js"))))
    mux.Handle("/", http.FileServer(http.Dir("./web/"))) 

    mux.HandleFunc("GET /{$}", handlers.Home)
	
	logger.Sugar().Info("Serving on http://%v ...", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Sugar().Fatal(err)
	}
}


//Параллелизм в Go Кэтрин Кокс-Будай Изучение Go Джона Боднера

