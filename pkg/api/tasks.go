package api

import (
	"net/http"
	"time"

	"github.com/len4ernova/go_final_project/pkg/db"
)

const (
	maxRows       = 50
	searchPattern = "02.01.2006"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler - вернуть список ближайших задач.
func (h *SrvHand) tasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []*db.Task
	h.Logger.Sugar().Info("START /api/tasks ", r.Method)
	search := r.URL.Query().Get("search")
	if len(search) > 0 {
		srchDate, err := verifySearchDate(search)
		if err != nil {
			// ищем по тексту
			tasks, err = db.TasksSearch(h.DB, maxRows, search, false)
		} else {
			// ищем по дате
			tasks, err = db.TasksSearch(h.DB, maxRows, srchDate, true)
		}
		if err != nil {
			// возвращает ошибку в JSON
			h.Logger.Sugar().Error(err)
			writeJson(w, reterror{Error: err.Error()})
			return
		}
		h.Logger.Sugar().Info("get", len(tasks), "tasks")
		writeJson(w, TasksResp{
			Tasks: tasks,
		})
		return
	}
	tasks, err := db.Tasks(h.DB, maxRows) // в параметре максимальное количество записей
	if err != nil {
		// возвращает ошибку в JSON
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

// verifySearchDate - проверка корректности введенной даты и конвертация в нужный формат.
func verifySearchDate(s string) (string, error) {
	t, err := time.Parse(searchPattern, s)
	if err != nil {
		return "", err
	}
	return t.Format(pattern), nil
}
