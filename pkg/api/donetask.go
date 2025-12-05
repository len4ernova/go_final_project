package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/len4ernova/go_final_project/pkg/db"
)

func (h *SrvHand) doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	fmt.Println("START doneTaskHandler", "strId = ", strId)
	id, err := strconv.Atoi(strId)
	if err != nil {
		fmt.Println(err)
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	task, err := db.GetTask(h.DB, id)
	if err != nil {
		fmt.Println(err)
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	fmt.Println(task)
	if len(task.Repeat) != 0 {
		err := db.DeleteTask(h.DB, id)
		if err != nil {
			writeJson(w, reterror{Error: err.Error()})
			return
		}
		writeJson(w, struct{}{})
		return
	}
	now := time.Now()
	nextdate, err := NextDate(now, task.Date, task.Repeat)
	fmt.Println("nextdate = ", nextdate, " task.Date, task.Repeat:", task.Date, task.Repeat)
	if err != nil {
		writeJson(w, struct{}{})
		return
	}
	task.Date = nextdate
	err = db.UpdateTask(h.DB, task)
	if err != nil {
		writeJson(w, struct{}{})
		return
	}
	writeJson(w, struct{}{})
}
