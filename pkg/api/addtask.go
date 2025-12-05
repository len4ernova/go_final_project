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
	fmt.Println("START /api/task", r.Method)
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

	err = checkDate(&task)
	if err != nil {
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	fmt.Println("11", task)
	idTask, err := db.AddTask(h.DB, &task)
	fmt.Println("11", "idTask: ", idTask)
	if err != nil {
		fmt.Println("22")
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	result := retid{
		Id: strconv.Itoa(int(idTask)),
	}
	fmt.Println("333")
	writeJson(w, result)
}

// writeJson - записать json.
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
	if len(task.Title) == 0 {
		return fmt.Errorf("value of <title> was empty")
	}
	now := time.Now()

	if len(task.Date) == 0 {
		task.Date = now.Format(pattern)
		fmt.Println("task.Date: ", task.Date)
	}
	fmt.Println("time - ", task.Date)
	t, err := time.Parse(pattern, task.Date)
	if err != nil {
		return err
	}
	fmt.Println("t", t)
	var next string
	if len(task.Repeat) != 0 {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	fmt.Println("next = ", next)
	fmt.Println("now: ", now, "task.Repeat", task.Repeat)
	fmt.Printf("afterNow(%v, %v) = %v\n", now, t, afterNow(now, t))
	if afterNow(now, t) {
		if task.Repeat == "" {

			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(pattern)
			fmt.Println("!!!", task.Date)
		} else {
			// в противном случае, берём вычисленную ранее следующую дату

			task.Date = next

			fmt.Println("!!!", task.Date)

		}
	}
	fmt.Println("@@@ Task ", task)

	return nil
}
