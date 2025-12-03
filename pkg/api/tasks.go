package api

import "net/http"

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func (h *SrvHand) tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		// здесь вызываете функцию, которая возвращает ошибку в JSON
		// её желательно было реализовать на предыдущем шаге
		// ...
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
