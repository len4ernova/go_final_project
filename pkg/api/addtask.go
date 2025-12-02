package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type req struct {
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func (h *SrvHand) taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println(r.Method)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		retErr(fmt.Sprintf("didn't get body request: %v", err))
		return
	}
	var b req
	json.Unmarshal([]byte(body), &b)
	fmt.Println(b)

	addTaskHandler()
}

func addTaskHandler() {

}

func retErr(s string) {
	w.Write()
}
