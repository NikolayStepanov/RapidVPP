package handlers

import (
	"net/http"

	"github.com/NikolayStepanov/RapidVPP/internal/delivery/http/handlers/Interfaces"
	"github.com/NikolayStepanov/RapidVPP/internal/delivery/http/handlers/acl"
	"github.com/NikolayStepanov/RapidVPP/internal/delivery/http/handlers/ip"
	"github.com/NikolayStepanov/RapidVPP/internal/delivery/http/handlers/vpp"
	"github.com/NikolayStepanov/RapidVPP/internal/service"
)

type Handler struct {
	router           *http.ServeMux
	vppHandler       *vpp.Handler
	interfaceHandler *interfaces.Handler
	ipHandler        *ip.Handler
	aclHandler       *acl.Handler
}

func NewHandler(info service.Info, inter service.Interface, IPServ service.IP, aclSer service.ACL) *Handler {
	handler := &Handler{
		router:           http.NewServeMux(),
		vppHandler:       vpp.NewHandler(info),
		interfaceHandler: interfaces.NewHandler(inter),
		ipHandler:        ip.NewHandler(IPServ),
		aclHandler:       acl.NewHandler(aclSer),
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
