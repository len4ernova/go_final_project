package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/len4ernova/go_final_project/pkg/db"
)

type reterror struct {
	Error string `json:"error"`
}
type retid struct {
	Id string `json:"id"`
}

func (h *SrvHand) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Sugar().Info("START /api/task ", r.Method)
	// Десериализация полученного запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Sugar().Errorf("didn't get body request: %v", err)
		writeJson(w, reterror{Error: fmt.Sprint("didn't get body request")}, http.StatusBadRequest)
		return
	}
	var task db.Task
	json.Unmarshal([]byte(body), &task)

	// Валидация данных
	err = checkDate(&task)
	if err != nil {
		h.Logger.Sugar().Errorf("task: %v, %v", task, err)
		writeJson(w, reterror{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	idTask, err := db.AddTask(h.DB, &task)
	if err != nil {
		h.Logger.Sugar().Error("task: %v, %v", task, err)
		writeJson(w, reterror{Error: "internal error"}, http.StatusInternalServerError)
		return
	}
	result := retid{
		Id: strconv.Itoa(int(idTask)),
	}
	h.Logger.Sugar().Info("/api/task: add task number ", result)
	writeJson(w, result, http.StatusOK)
}

// writeJson - записать json.
func writeJson(w http.ResponseWriter, data any, statusCode int) {
	jsondata, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "error encoding response", statusCode)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsondata)
	if err != nil {
		http.Error(w, "error writing data", http.StatusInternalServerError)
	}
}

// checkDate - проверить на корректность полученное значение task.Date.
func checkDate(task *db.Task) error {
	if len(task.Title) == 0 {
		return fmt.Errorf("value of <title> was empty")
	}
	now := time.Now()

	if len(task.Date) == 0 {
		task.Date = now.Format(pattern)
	}
	t, err := time.Parse(pattern, task.Date)
	if err != nil {
		return err
	}

	var next string
	if len(task.Repeat) != 0 {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(pattern)
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}
	return nil
}
