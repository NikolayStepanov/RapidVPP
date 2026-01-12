package interfaces

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	interfaces, err := h.inter.List(r.Context())
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

func (h *Handler) CreateLoopback(w http.ResponseWriter, r *http.Request) {
	createLoopbackResponse, err := h.inter.CreateLoopback(r.Context())
	if err != nil {
		logger.Error("Failed to create loopback interfaces", zap.Error(err))
		http.Error(w, "Failed to create loopback interfaces", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createLoopbackResponse); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SetInterfaceState(w http.ResponseWriter, r *http.Request) {
	var req SetInterfaceStateRequests
	idStr := r.PathValue("id")
	ifIndex, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warn("Invalid index interface in request", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Invalid index interface", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	err = h.inter.SetInterfaceAdminState(r.Context(), uint32(ifIndex), req.AdminUp)
	if err != nil {
		logger.Error("Failed to set interface state", zap.Error(err))
		http.Error(w, "Failed to set interface state", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
