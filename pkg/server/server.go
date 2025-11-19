package server

import (
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type Settings struct {
	Ip   string
	Port int
}

const webDir = "./web"
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

	fileServer := http.FileServer(http.Dir("./web"))
	
	// Указываем обработчики
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	}) //index.html
	 
    mux.Handle("GET /js/scripts.min.js/", fileServer)
    mux.Handle("GET /css/style.css/", fileServer)
    mux.Handle("GET /favicon.ico/", fileServer)

	// mux.HandleFunc("GET /css/style.css", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./web/css/style.css")
	// })
	// mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./web/favicon.ico")
	// })
	// mux.HandleFunc("/web/js/scripts.min.js", handlers.js)
	// mux.HandleFunc("/web/css/style.css", handlers.css)
	// mux.HandleFunc("/web/favicon.ico", handlers.css)
	
	logger.Sugar().Info("Serving on http://%v ...", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Sugar().Fatal(err)
	}
}
