package handlers

import (
	"net/http"

	"github.com/NikolayStepanov/RapidVPP/internal/delivery/http/handlers/vpp"
	"github.com/NikolayStepanov/RapidVPP/internal/service"
)

type Handler struct {
	router     *http.ServeMux
	vppHandler *vpp.Handler
}

func NewHandler(info service.Info) *Handler {
	handler := &Handler{
		router:     http.NewServeMux(),
		vppHandler: vpp.NewHandler(info),
	}

	handler.setupRoutes()
	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) GetRouter() *http.ServeMux {
	return h.router
}
