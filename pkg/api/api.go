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

func (h *SrvHand) Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", h.nextDayHandler)       // рассчитать следующую дату
	mux.HandleFunc("POST /api/task", h.addTaskHandler)      // добавляет задачу в формате JSON
	mux.HandleFunc("GET /api/tasks", h.tasksHandler)        // возвращать список ближайших задач в формате JSON
	mux.HandleFunc("GET /api/task", h.getTaskHandler)       // возвращает задачу по ID
	mux.HandleFunc("PUT /api/task", h.putTaskHandler)       // изменяет задачу по ID
	mux.HandleFunc("DELETE /api/task", h.deleteTaskHandler) // удалить задачу по ID
	mux.HandleFunc("PUT /api/task/done", h.doneTaskHandler) // задача выполнена

}
