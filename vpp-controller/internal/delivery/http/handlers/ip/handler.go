package ip

import (
	"encoding/json"
	"net/http"

	"github.com/NikolayStepanov/RapidVPP/internal/service"
	"github.com/NikolayStepanov/RapidVPP/pkg/logger"
	"go.uber.org/zap"
)

type Handler struct {
	ip service.IP
}

func NewHandler(ip service.IP) *Handler {
	return &Handler{ip: ip}
}

func (h *Handler) AddRoute(w http.ResponseWriter, r *http.Request) {
	var req AddRouteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	route, err := req.ToDomain()
	if err != nil {
		logger.Error("Failed to add route", zap.Error(err))
		http.Error(w, "Failed to add route", http.StatusBadRequest)
		return
	}
	err = h.ip.AddRoute(r.Context(), route)
	if err != nil {
		logger.Error("Failed to add route", zap.Error(err))
		http.Error(w, "Failed to add route", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteRoute(w http.ResponseWriter, r *http.Request) {
	var req AddRouteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	route, err := req.ToDomain()
	if err != nil {
		logger.Error("Failed to delete route", zap.Error(err))
		http.Error(w, "Failed to delete route", http.StatusBadRequest)
		return
	}
	err = h.ip.DeleteRoute(r.Context(), route)
	if err != nil {
		logger.Error("Failed to delete route", zap.Error(err))
		http.Error(w, "Failed to delete route", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
