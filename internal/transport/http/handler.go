package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Router *mux.Router
}

type Response struct {
	Message string
	Error   string
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SetUpRoutes() {
	log.Info("Setting up routes")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(Response{Message: "Healthy"}); err != nil {
			log.Panic(err)
		}
	})
}
