package api

import (
	"net/http"
	"strconv"

	"github.com/len4ernova/go_final_project/pkg/db"
)

func (h *SrvHand) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	task, err := db.GetTask(h.db, id)
}
