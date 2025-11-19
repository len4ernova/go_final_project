package handlers

import (
	"fmt"
	"net/http"
)

func GetIndexHTML(w http.ResponseWriter, r *http.Request) {
	// data, err := os.ReadFile("web/index.html")
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("ошибка при считывании файла: %v", err), http.StatusInternalServerError)
	// 	return
	// }
	fmt.Println("------")
	http.ServeFile(w, r, "./web/index.html")
	// http.ServeFile(w, r, "/web/js/scripts.min.js")
	// http.ServeFile(w, r, "/web/css/style.css")
	// http.ServeFile(w, r, "./web/favicon.ico")
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.WriteHeader(http.StatusOK)
	// w.Write(data)
}
