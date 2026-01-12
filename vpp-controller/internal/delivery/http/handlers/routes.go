package handlers

func (h *Handler) setupRoutes() {
	h.router.HandleFunc("GET /vpp/version", h.vppHandler.GetVersion)
	h.router.HandleFunc("GET /interfaces/", h.interfaceHandler.List)
	h.router.HandleFunc("POST /interfaces/loopback", h.interfaceHandler.CreateLoopback)
}
