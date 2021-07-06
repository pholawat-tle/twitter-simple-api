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
	h.Router.HandleFunc("/api/tweet", h.CreateTweet).Methods("POST")
	h.Router.HandleFunc("/api/tweet/{id}", h.GetTweetByID).Methods("GET")
	h.Router.HandleFunc("/api/tweet/{id}", h.DeleteTweetByID).Methods("DELETE")

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
	rw.WriteHeader(http.StatusOK)
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
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(tweet); err != nil {
		log.Panic(err)
	}
}

func (h *Handler) CreateTweet(rw http.ResponseWriter, r *http.Request) {
	var tweet tweet.Tweet
	if err := json.NewDecoder(r.Body).Decode(&tweet); err != nil {
		SendErrorResponse(rw, "Failed to parse tweet from body", err)
		return
	}

	tweet.Likes = 0
	tweet.Author = GetIP(r)

	tweet, err := h.Service.CreateTweet(tweet)
	if err != nil {
		SendErrorResponse(rw, "Failed to create tweet", err)
		return
	}
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(Response{Message: "Your tweet has been posted"}); err != nil {
		log.Panic(err)
	}
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func (h *Handler) DeleteTweetByID(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	uint_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		SendErrorResponse(rw, "Failed to parse ID to uint", err)
		return
	}

	err = h.Service.DeleteTweetByID(uint(uint_id))
	if err != nil {
		SendErrorResponse(rw, "Failed to retrieve tweet by ID", err)
		return
	}
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(Response{Message: "The tweet has been deleted"}); err != nil {
		log.Panic(err)
	}
}

func SendErrorResponse(rw http.ResponseWriter, msg string, err error) {
	log.Error(err.Error())
	rw.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(rw).Encode(Response{Message: msg, Error: err.Error()})
}
