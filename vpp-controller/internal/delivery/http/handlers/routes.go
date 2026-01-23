package handlers

func (h *Handler) setupRoutes() {
	h.router.HandleFunc("GET /vpp/version", h.vppHandler.GetVersion)
	h.router.HandleFunc("GET /interfaces/", h.interfaceHandler.List)
	h.router.HandleFunc("POST /interfaces/loopback", h.interfaceHandler.CreateLoopback)
	h.router.HandleFunc("POST /interfaces/{id}/state", h.interfaceHandler.SetInterfaceState)
	h.router.HandleFunc("DELETE /interfaces/{id}", h.interfaceHandler.DeleteLoopback)
	h.router.HandleFunc("POST /interfaces/{id}/ip", h.interfaceHandler.AddInterfaceIP)
	h.router.HandleFunc("POST /routes", h.ipHandler.AddRoute)
	h.router.HandleFunc("DELETE /routes", h.ipHandler.DeleteRoute)
	h.router.HandleFunc("GET /routes", h.ipHandler.List)
	h.router.HandleFunc("GET /routes/{vrf}", h.ipHandler.Get)
	h.router.HandleFunc("POST /vrf", h.ipHandler.CreateVRF)
	h.router.HandleFunc("GET /vrf", h.ipHandler.ListVRF)
	h.router.HandleFunc("DELETE /vrf/{id}", h.ipHandler.DeleteVRF)
	h.router.HandleFunc("POST /acl", h.aclHandler.Create)
	h.router.HandleFunc("PUT /acl/{id}", h.aclHandler.Update)
}
