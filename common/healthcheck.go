package common

import (
	"net/http"
)

func HandleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HEALTHY"))
}

type HealthcheckServer struct{}

func NewHealthcheckServer() *HealthcheckServer {
	return &HealthcheckServer{}
}

func (h *HealthcheckServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleHealthcheck(w, r)
}
