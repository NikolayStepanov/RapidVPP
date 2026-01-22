package ip

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
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
		http.Error(w, "Failed to add route", http.StatusInternalServerError)
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
		http.Error(w, "Failed to delete route", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	vrfStr := r.URL.Query().Get("vrf")
	var vrf uint32 = 0

	if vrfStr != "" {
		vrfInt, err := strconv.ParseUint(vrfStr, 10, 32)
		if err != nil {
			logger.Error("Invalid VRF parameter", zap.String("vrf", vrfStr), zap.Error(err))
			http.Error(w, "Invalid VRF parameter. Must be a number", http.StatusBadRequest)
			return
		}
		vrf = uint32(vrfInt)
	}
	routes, err := h.ip.ListRoutes(r.Context(), vrf)
	if err != nil {
		logger.Error("Failed to get list routes", zap.Error(err))
		http.Error(w, "Failed to get list routes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(routes); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vrfStr := r.PathValue("vrf")
	if vrfStr == "" {
		logger.Error("VRF not found in path", zap.String("vrf", vrfStr))
		http.Error(w, "VRF not found in path", http.StatusBadRequest)
		return
	}

	vrfInt, err := strconv.ParseUint(vrfStr, 10, 32)
	if err != nil {
		logger.Error("Invalid VRF parameter", zap.String("vrf", vrfStr), zap.Error(err))
		http.Error(w, "Invalid VRF ID", http.StatusBadRequest)
		return
	}
	vrf := uint32(vrfInt)

	prefixStr := r.URL.Query().Get("prefix")
	if prefixStr == "" {
		logger.Error("Invalid prefix parameter", zap.String("prefix", prefixStr))
		http.Error(w, "Invalid prefix parameter", http.StatusBadRequest)
		return
	}

	ip, netw, err := net.ParseCIDR(prefixStr)
	if err != nil {
		logger.Error("Invalid prefix query param", zap.String("prefix", prefixStr), zap.Error(err))
		http.Error(w, "Invalid prefix format. Use: address/prefix", http.StatusBadRequest)
		return
	}

	ones, _ := netw.Mask.Size()

	dst := domain.IPWithPrefix{
		Address: ip.String(),
		Prefix:  uint8(ones),
	}

	route, err := h.ip.GetRoute(r.Context(), dst, vrf)
	if err != nil {
		logger.Error("Failed to get route",
			zap.String("prefix", prefixStr),
			zap.Uint32("vrf", vrf),
			zap.Error(err))
		http.Error(w, "Route not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(route); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateVRF(w http.ResponseWriter, r *http.Request) {
	var req CreateVRFRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	err := h.ip.CreateVRF(r.Context(), req.Id, req.Name)
	if err != nil {
		logger.Error("Failed to create VRF", zap.Error(err))
		http.Error(w, "Failed to create VRF", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListVRF(w http.ResponseWriter, r *http.Request) {
	vrfs, err := h.ip.ListVRF(r.Context())
	if err != nil {
		logger.Error("Failed to list VRF", zap.Error(err))
		http.Error(w, "Failed to list VRF", http.StatusInternalServerError)
		return
	}

	response := VRFToResponse(vrfs)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteVRF(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warn("Invalid ID VRF in request", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Invalid ID VRF", http.StatusBadRequest)
		return
	}

	if err := h.ip.DeleteVRF(r.Context(), uint32(id)); err != nil {
		logger.Error("Failed to delete VRF", zap.Uint64("ID", id), zap.Error(err))
	} else {
		logger.Info("VRF deleted successfully", zap.Uint64("ID", id))
	}

	w.WriteHeader(http.StatusNoContent)
}
