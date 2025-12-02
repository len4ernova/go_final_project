package api

import (
	"net/http"

	"go.uber.org/zap"
)

type SrvHand struct {
	Logger *zap.Logger
}

const pattern = "20060102"

func (h *SrvHand) Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", h.nextDayHandler)
	mux.HandleFunc("POST /api/task", h.taskHandler)

}
