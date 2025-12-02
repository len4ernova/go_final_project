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

// func (h *SrvHand) taskHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		fmt.Println(r.Method)
// 	}
// 	fmt.Println(r.Method)

// }

type reterror struct {
	Error string `json:"error"`
}
type retid struct {
	Id string `json:"id"`
}

func (h *SrvHand) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Десериализация полученного запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, reterror{Error: fmt.Sprintf("didn't get body request: %v", err)})
		return
	}
	var task db.Task
	json.Unmarshal([]byte(body), &task)
	fmt.Println("01", task)

	// Валидация данных
	if len(task.Title) == 0 {
		writeJson(w, reterror{Error: "value of <title> was empty"})
		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	fmt.Println("11", task)
	idTask, err := db.AddTask(h.DB, &task)
	fmt.Println("11", "idTask: ", idTask)
	if err != nil {
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	result := retid{
		Id: strconv.Itoa(int(idTask)),
	}
	writeJson(w, result)

}

// retErr - возврат ошибки.
func writeJson(w http.ResponseWriter, data any) {
	jsondata, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsondata)

}

// checkDate - проверить на корректность полученное значение task.Date.
func checkDate(task *db.Task) error {
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

	if !afterNow(now, t) {
		if task.Repeat == "" {

			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(pattern)
			fmt.Println("!!!", task)
		} else {
			// в противном случае, берём вычисленную ранее следующую дату

			task.Date = next
			fmt.Println("!!!", task)

		}
	}

	return nil
}
