package rest

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/DexScen/VKtestTask/backend/internal/domain"
	"github.com/gorilla/mux"
)

type Containers interface {
	GetContainers(ctx context.Context, list *domain.ListContainer) error
	PostContainers(ctx context.Context, list *domain.ListContainer) error
}

type Handler struct {
	containersService Containers
}

func NewHandler(containers Containers) *Handler {
	return &Handler{
		containersService: containers,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(loggingMiddleware)

	links := r.PathPrefix("/containers").Subrouter()
	{
		links.HandleFunc("", h.GetContainers).Methods(http.MethodGet)
		links.HandleFunc("", h.PostContainers).Methods(http.MethodPost)
	}

	return r
}

func (h *Handler) GetContainers(w http.ResponseWriter, r *http.Request) {
	var list domain.ListContainer
	if err := h.containersService.GetContainers(context.TODO(), &list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("getContainers error:", err)
		return
	}

	if jsonResp, err := json.Marshal(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("getContainers error:", err)
		return
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Write(jsonResp)
	}
}

func (h *Handler) PostContainers(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("postContainers error:", err)
		return
	}
	defer r.Body.Close()

	var list domain.ListContainer
	if err = json.Unmarshal(reqBytes, &list); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("postContainers error:", err)
		return
	}

	if err := h.containersService.PostContainers(context.TODO(), &list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("postContainers error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
