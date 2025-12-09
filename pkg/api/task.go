package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/len4ernova/go_final_project/pkg/db"
)

// getTaskHandler - вернуть задачу по ID.
func (h *SrvHand) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	h.Logger.Sugar().Info("START: /api/task", strId, r.Method)
	id, err := strconv.Atoi(strId)
	if err != nil {
		h.Logger.Sugar().Error(err)
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	task, err := db.GetTask(h.DB, id)
	if err != nil {
		h.Logger.Sugar().Error(err)
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	writeJson(w, task)
}

// putTaskHandler - изменить задачу по ID.
func (h *SrvHand) putTaskHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Sugar().Info("START /api/task", r.Method)
	// Десериализация полученного запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Sugar().Error(fmt.Sprintf("didn't get body request: %v", err))
		writeJson(w, reterror{Error: fmt.Sprintf("didn't get body request: %v", err)})
		return
	}
	var task db.Task
	json.Unmarshal([]byte(body), &task)

	// Валидация данных
	if len(task.ID) == 0 {
		h.Logger.Sugar().Error("value of <title> was empty")
		writeJson(w, reterror{Error: "value of <title> was empty"})
		return
	}
	err = checkDate(&task)
	if err != nil {
		h.Logger.Sugar().Error(err)
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	err = db.UpdateTask(h.DB, &task)
	if err != nil {
		h.Logger.Sugar().Error(err)
		writeJson(w, reterror{Error: err.Error()})
		return
	}

	h.Logger.Sugar().Info("task changed ok")
	writeJson(w, struct{}{})
}
