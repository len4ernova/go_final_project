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

	mux.HandleFunc("GET /", handlers.GetIndexHTML) //index.html
	/*     http://localhost:7540/js/scripts.min.js возвращает ./web/js/scripts.min.js;
	http://localhost:7540/css/style.css возвращает ./web/css/style.css;
	http://localhost:7540/favicon.ico возвращает ./web/favicon.ico.
	*/
	logger.Sugar().Info("Serving on http://%v ...", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Sugar().Fatal(err)
	}
}
