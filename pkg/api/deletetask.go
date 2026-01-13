package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/len4ernova/go_final_project/pkg/db"
)

// deleteTaskHandler - удалить задачу по ID.
func (h *SrvHand) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	h.Logger.Sugar().Infof("START /api/task %v strId = %v", r.Method, strId)
	if len(strId) == 0 {
		h.Logger.Sugar().Error("empty value <id>")
		writeJson(w, reterror{Error: "empty value <id>"}, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		h.Logger.Sugar().Error(err)
		writeJson(w, reterror{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	err = db.DeleteTask(h.DB, id)
	if err != nil {
		h.Logger.Sugar().Error("error deletion: ", err)
		if errors.Is(err, db.ErrDeleteZeroRows) {
			writeJson(w, reterror{Error: "rows didn't found"}, http.StatusNotFound)
			return
		}
		writeJson(w, reterror{Error: "task deletion bug"}, http.StatusInternalServerError)
		return
	}

	h.Logger.Sugar().Info(strId, " deletion ok")
	writeJson(w, struct{}{}, http.StatusOK)
}
