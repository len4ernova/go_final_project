package api

import (
	"net/http"
	"strconv"

	"github.com/len4ernova/go_final_project/pkg/db"
)

// deleteTaskHandler - удалить задачу по ID.
func (h *SrvHand) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	h.Logger.Sugar().Info("START /api/task", "strId = ", strId, r.Method)
	if len(strId) == 0 {
		h.Logger.Sugar().Error("empty value <id>")
		writeJson(w, reterror{Error: "empty value <id>"})
		return
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		h.Logger.Sugar().Error(err)
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	err = db.DeleteTask(h.DB, id)
	if err != nil {
		h.Logger.Sugar().Error(err)
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	h.Logger.Sugar().Info(strId, "deleted ok")
	writeJson(w, struct{}{})
}
