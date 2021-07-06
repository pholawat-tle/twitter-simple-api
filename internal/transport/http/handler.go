package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"twitter/internal/tweet"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Router  *mux.Router
	Service *tweet.Service
}

type Response struct {
	Message string
	Error   string
}

func NewHandler(s *tweet.Service) *Handler {
	return &Handler{Service: s}
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
			},
		).Info("Request Received")
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(rw, r)
	})
}

func (h *Handler) SetUpRoutes() {
	log.Info("Setting up routes")
	h.Router = mux.NewRouter()
	h.Router.Use(LogMiddleware)

	h.Router.HandleFunc("/api/tweet", h.GetAllTweet).Methods("GET")
	h.Router.HandleFunc("/api/tweet/{id}", h.GetTweetByID).Methods("GET")

	h.Router.HandleFunc("/api/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(rw).Encode(Response{Message: "Healthy"}); err != nil {
			log.Panic(err)
		}
	})
}

func (h *Handler) GetAllTweet(rw http.ResponseWriter, r *http.Request) {
	tweets, err := h.Service.GetAllTweet()
	if err != nil {
		SendErrorResponse(rw, "Failed to retrieve all tweets", err)
		return
	}
	if err := json.NewEncoder(rw).Encode(tweets); err != nil {
		log.Panic(err)
	}
}

func (h *Handler) GetTweetByID(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	uint_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		SendErrorResponse(rw, "Failed to parse ID to uint", err)
		return
	}

	tweet, err := h.Service.GetTweetByID(uint(uint_id))
	if err != nil {
		SendErrorResponse(rw, "Failed to retrieve tweet by ID", err)
		return
	}

	if err := json.NewEncoder(rw).Encode(tweet); err != nil {
		log.Panic(err)
	}
}

func SendErrorResponse(rw http.ResponseWriter, msg string, err error) {
	log.Error(err.Error())
	rw.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(rw).Encode(Response{Message: msg, Error: err.Error()})
}
