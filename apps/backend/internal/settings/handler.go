package settings

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	slog.Info("request received", "method", r.Method, "path", r.URL.Path)
	s, err := h.svc.Get(r.Context())
	if err != nil {
		h.writeError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("request completed", "method", r.Method, "path", r.URL.Path, "status", http.StatusOK)
	json.NewEncoder(w).Encode(s)
}

func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	slog.Info("request received", "method", r.Method, "path", r.URL.Path)
	var s Settings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		h.writeError(w, "VALIDATION_ERROR", err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.Update(r.Context(), &s); err != nil {
		h.writeError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("request completed", "method", r.Method, "path", r.URL.Path, "status", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) writeError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
		"correlationId": uuid.New().String(),
	}

	json.NewEncoder(w).Encode(resp)
}