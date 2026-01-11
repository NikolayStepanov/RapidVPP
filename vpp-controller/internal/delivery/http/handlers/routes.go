package handlers

func (h *Handler) setupRoutes() {
	h.router.HandleFunc("GET /vpp/version", h.vppHandler.GetVersion)
	h.router.HandleFunc("GET /interface/list", h.interfaceHandler.List)
}
