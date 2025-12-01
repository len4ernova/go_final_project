package api

import (
	"net/http"
	"text/template"
)

const indexHTML = "./web/index.html"

func (h *SrvHand) Index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles(indexHTML)
	if err != nil {
		h.Logger.Sugar().Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Then we use the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		h.Logger.Sugar().Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
