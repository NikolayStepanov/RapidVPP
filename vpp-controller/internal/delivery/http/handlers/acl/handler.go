package acl

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"github.com/NikolayStepanov/RapidVPP/internal/service"
	"github.com/NikolayStepanov/RapidVPP/pkg/logger"
	"go.uber.org/zap"
)

type Handler struct {
	acl service.ACL
}

func NewHandler(acl service.ACL) *Handler {
	return &Handler{acl: acl}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	aclRules, err := ConvertRulesRequestToDomain(req.Rules)
	if err != nil {
		logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	AclID, err := h.acl.Create(r.Context(), req.Name, aclRules)
	if err != nil {
		logger.Error("Failed to create acl", zap.Error(err))
		http.Error(w, "Failed to create acl", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(AclID); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req UpdateRequest
	idStr := r.PathValue("id")
	aclID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warn("Invalid index acl in request", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Invalid index acl", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	aclRules, err := ConvertRulesRequestToDomain(req.Rules)
	if err != nil {
		logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	id := domain.AclID(uint32(aclID))
	err = h.acl.Update(r.Context(), id, aclRules)
	if err != nil {
		logger.Error("Failed to update acl", zap.Error(err))
		http.Error(w, "Failed to update acl", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	aclID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warn("Invalid index acl in request", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Invalid index acl", http.StatusBadRequest)
		return
	}
	id := domain.AclID(uint32(aclID))
	if err := h.acl.Delete(r.Context(), id); err != nil {
		logger.Error("Failed to delete ACL", zap.Uint64("aclID", aclID), zap.Error(err))
	} else {
		logger.Info("ACL deleted successfully", zap.Uint64("aclID", aclID))
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var resp ListACLResponse
	acls, err := h.acl.List(r.Context())
	if err != nil {
		logger.Error("Failed to list ACL", zap.Error(err))
		http.Error(w, "Failed to list ACL", http.StatusInternalServerError)
		return
	}
	resp = ListACLResponse{InfosToResponse(acls)}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
