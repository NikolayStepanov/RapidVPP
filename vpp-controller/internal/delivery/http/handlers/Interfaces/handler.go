package Interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/NikolayStepanov/RapidVPP/internal/service"
	"github.com/NikolayStepanov/RapidVPP/pkg/logger"
	"go.uber.org/zap"
)

type Handler struct {
	inter service.Interface
}

func NewHandler(inter service.Interface) *Handler {
	return &Handler{inter: inter}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	interfaces, err := h.inter.List()
	if err != nil {
		logger.Error("Failed to get list interfaces", zap.Error(err))
		http.Error(w, "Failed to get list interfaces", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(interfaces); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
