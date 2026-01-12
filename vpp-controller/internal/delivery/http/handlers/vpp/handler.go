package vpp

import (
	"encoding/json"
	"net/http"

	"github.com/NikolayStepanov/RapidVPP/internal/service"
	"github.com/NikolayStepanov/RapidVPP/pkg/logger"
	"go.uber.org/zap"
)

type Handler struct {
	infoService service.Info
}

func NewHandler(info service.Info) *Handler {
	return &Handler{infoService: info}
}

func (h *Handler) GetVersion(w http.ResponseWriter, r *http.Request) {
	version, err := h.infoService.GetVersion(r.Context())
	if err != nil {
		logger.Error("Failed to get version", zap.Error(err))
		http.Error(w, "Failed to get version", http.StatusBadRequest)
		return
	}

	response := ToVersionResponse(version)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
