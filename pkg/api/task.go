package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/len4ernova/go_final_project/pkg/db"
)

func (h *SrvHand) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	fmt.Println("START: getTaskHandler", strId)
	id, err := strconv.Atoi(strId)
	if err != nil {
		writeJson(w, reterror{Error: err.Error()})
	}
	task, err := db.GetTask(h.DB, id)
	writeJson(w, task)
}

func (h *SrvHand) putTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Десериализация полученного запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, reterror{Error: fmt.Sprintf("didn't get body request: %v", err)})
		return
	}
	var task db.Task
	json.Unmarshal([]byte(body), &task)

	// Валидация данных
	if len(task.ID) == 0 {
		writeJson(w, reterror{Error: "value of <title> was empty"})
		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	err = db.UpdateTask(h.DB, &task)
	if err != nil {
		writeJson(w, reterror{Error: err.Error()})
		return
	}

	writeJson(w, struct{}{})
}
