package source

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	src := &Source{URL: req.URL}
	if err := h.service.Create(r.Context(), src); err != nil {
		if err.Error() == "Duplicate detected" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		// Log the actual error for debugging
		fmt.Printf("Error creating source: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(src)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	sources, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Ensure we return [] instead of null for empty list
	if sources == nil {
		sources = []Source{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sources)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ReSync(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.ReSync(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}