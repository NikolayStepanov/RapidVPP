package handlers

func (h *Handler) setupRoutes() {
	h.router.HandleFunc("GET /vpp/version", h.vppHandler.GetVersion)
	h.router.HandleFunc("GET /interfaces/", h.interfaceHandler.List)
	h.router.HandleFunc("POST /interfaces/loopback", h.interfaceHandler.CreateLoopback)
	h.router.HandleFunc("POST /interfaces/{id}/state", h.interfaceHandler.SetInterfaceState)
	h.router.HandleFunc("DELETE /interfaces/{id}", h.interfaceHandler.DeleteLoopback)
	h.router.HandleFunc("POST /interfaces/{id}/ip", h.interfaceHandler.AddInterfaceIP)
}
