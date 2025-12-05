package api

import (
	"database/sql"
	"net/http"

	"go.uber.org/zap"
)

type SrvHand struct {
	Logger *zap.Logger
	DB     *sql.DB
}

const pattern = "20060102"

// Init - конечные точки и вызов обработчиков.
func (h *SrvHand) Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", h.nextDayHandler)        // рассчитать следующую дату
	mux.HandleFunc("POST /api/task", h.addTaskHandler)       // добавить задачу в формате JSON
	mux.HandleFunc("GET /api/tasks", h.tasksHandler)         // вернуть список ближайших задач в формате JSON
	mux.HandleFunc("GET /api/task", h.getTaskHandler)        // вернуть задачу по ID
	mux.HandleFunc("PUT /api/task", h.putTaskHandler)        // изменить задачу по ID
	mux.HandleFunc("DELETE /api/task", h.deleteTaskHandler)  // удалить задачу по ID
	mux.HandleFunc("POST /api/task/done", h.doneTaskHandler) // задача выполнена

}
