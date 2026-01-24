package interfaces

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"github.com/NikolayStepanov/RapidVPP/internal/service"
	"github.com/NikolayStepanov/RapidVPP/internal/service/vpp/interfaces"
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
		http.Error(w, "Failed to create loopback interfaces", http.StatusInternalServerError)
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
		http.Error(w, "Failed to set interface state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteLoopback(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	ifIndex, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warn("Invalid index interface in request", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Invalid index interface", http.StatusBadRequest)
		return
	}

	if err := h.inter.DeleteLoopback(r.Context(), uint32(ifIndex)); err != nil {
		logger.Error("Failed to delete loopback interface", zap.Uint64("ifIndex", ifIndex), zap.Error(err))
	} else {
		logger.Info("Loopback interface deleted successfully", zap.Uint64("ifIndex", ifIndex))
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AddInterfaceIP(w http.ResponseWriter, r *http.Request) {
	var req AddIPRequest
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
	IPWithPrefix := domain.IPWithPrefix{req.Address, req.Prefix}
	err = h.inter.SetInterfaceIP(r.Context(), uint32(ifIndex), IPWithPrefix)
	switch {
	case errors.Is(err, interfaces.ErrNotFound):
		logger.Error("Failed interface not found", zap.Error(err))
		http.Error(w, "interface not found", http.StatusNotFound)
	case errors.Is(err, interfaces.ErrAlreadyExists):
		logger.Error("Failed IP address already exists", zap.Error(err))
		http.Error(w, "IP address already exists", http.StatusConflict)
	case err != nil:
		logger.Error("Failed to add interface IP", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
	if err != nil {
		logger.Error("Failed to set interface state", zap.Error(err))
		http.Error(w, "Failed to set interface state", http.StatusBadRequest)
		return
	}
}

func (h *Handler) AttachACL(w http.ResponseWriter, r *http.Request) {
	var req AttachACLRequest
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

	err = h.inter.AttachACL(r.Context(), uint32(ifIndex), req.AclId, req.Direction)
	if err != nil {
		logger.Error("Failed to attach ACL interface", zap.Error(err))
		http.Error(w, "Failed to attach ACL interface", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DetachACL(w http.ResponseWriter, r *http.Request) {
	var req DetachACLRequest

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

	err = h.inter.DetachACL(r.Context(), uint32(ifIndex), req.AclId, req.Direction)
	if err != nil {
		logger.Error("Failed to detach ACL interface", zap.Error(err))
		http.Error(w, "Failed to detach ACL interface", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListACL(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	ifIndex, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warn("Invalid index interface in request", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Invalid index interface", http.StatusBadRequest)
		return
	}

	aclList, err := h.inter.ListACL(r.Context(), uint32(ifIndex))
	if err != nil {
		logger.Error("Failed to get list acl interface", zap.Error(err))
		http.Error(w, "Failed to get list acl interface", http.StatusBadRequest)
		return
	}
	aclListDTO := ACLInterfaceListToDTO(aclList)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(aclListDTO); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
