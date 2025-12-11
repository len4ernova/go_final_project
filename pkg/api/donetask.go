package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/len4ernova/go_final_project/pkg/db"
)

// doneTaskHandler - задача выполнена.
func (h *SrvHand) doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	h.Logger.Sugar().Info("START /api/task/done ", r.Method, " strId = ", strId)
	id, err := strconv.Atoi(strId)
	if err != nil {
		h.Logger.Sugar().Errorf("<id>= %v : %v", strId, err)
		writeJson(w, reterror{Error: "bug in <id>"}, http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(h.DB, id)
	if err != nil {
		h.Logger.Sugar().Errorf("error of data query: %v", err)
		writeJson(w, reterror{Error: "internal error"}, http.StatusInternalServerError)
		return
	}
	if len(task.Repeat) == 0 {
		err := db.DeleteTask(h.DB, id)
		if err != nil {
			h.Logger.Sugar().Errorf("task deletion bug: %v", err)
			writeJson(w, reterror{Error: "internal error"}, http.StatusInternalServerError)
			return
		}
		writeJson(w, struct{}{}, http.StatusOK)
		return
	}
	now := time.Now()
	nextdate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		h.Logger.Sugar().Error("doneTaskHandler: NextDate", err)
		writeJson(w, struct{}{}, http.StatusInternalServerError)
		return
	}
	task.Date = nextdate
	err = db.UpdateTask(h.DB, task)
	if err != nil {
		h.Logger.Sugar().Error("doneTaskHandler: UpdateTask", err) // TODO
		writeJson(w, struct{}{}, http.StatusInternalServerError)
		return
	}
	h.Logger.Sugar().Info("task ", strId, " changed")
	writeJson(w, struct{}{}, http.StatusOK)
}
