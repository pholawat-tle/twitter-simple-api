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

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
			},
		).Info("Request Received")
		next.ServeHTTP(rw, r)
	})
}

func (h *Handler) SetUpRoutes() {
	log.Info("Setting up routes")
	h.Router = mux.NewRouter()
	h.Router.Use(LogMiddleware)

	h.Router.HandleFunc("/api/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(rw).Encode(Response{Message: "Healthy"}); err != nil {
			log.Panic(err)
		}
	})
}
