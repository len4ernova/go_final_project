package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/len4ernova/go_final_project/pkg/db"
)

func (h *SrvHand) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	fmt.Println("START deleteTaskHandler", "strId = ", strId)
	id, err := strconv.Atoi(strId)
	if err != nil {
		fmt.Println(err)
		//writeJson(w, reterror{Error: err.Error()})
		writeJson(w, struct{}{})
		return
	}
	err = db.DeleteTask(h.DB, id)
	if err != nil {
		fmt.Println(err)
		writeJson(w, reterror{Error: err.Error()})
		return
	}
	fmt.Println(struct{}{})
	writeJson(w, struct{}{})
}
