package api

import (
	"net/http"

	"github.com/len4ernova/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func (h *SrvHand) tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(h.DB, 50) // в параметре максимальное количество записей
	if err != nil {
		// возвращает ошибку в JSON
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
