package handlers

func (h *Handler) setupRoutes() {
	h.router.HandleFunc("GET /vpp/version", h.vppHandler.GetVersion)

	h.router.HandleFunc("GET /interfaces/", h.interfaceHandler.List)
	h.router.HandleFunc("POST /interfaces/loopback", h.interfaceHandler.CreateLoopback)
	h.router.HandleFunc("POST /interfaces/{id}/state", h.interfaceHandler.SetInterfaceState)
	h.router.HandleFunc("DELETE /interfaces/{id}", h.interfaceHandler.DeleteLoopback)
	h.router.HandleFunc("POST /interfaces/{id}/ip", h.interfaceHandler.AddInterfaceIP)

	h.router.HandleFunc("POST /interfaces/{id}/acl", h.interfaceHandler.AttachACL)
	//h.router.HandleFunc("DELETE /interfaces/{id}/acl/{acl}",)
	//h.router.HandleFunc("GET /interfaces/{id}/acl",)

	h.router.HandleFunc("POST /routes", h.ipHandler.AddRoute)
	h.router.HandleFunc("DELETE /routes", h.ipHandler.DeleteRoute)
	h.router.HandleFunc("GET /routes", h.ipHandler.List)
	h.router.HandleFunc("GET /routes/{vrf}", h.ipHandler.Get)

	h.router.HandleFunc("POST /vrf", h.ipHandler.CreateVRF)
	h.router.HandleFunc("GET /vrf", h.ipHandler.ListVRF)
	h.router.HandleFunc("DELETE /vrf/{id}", h.ipHandler.DeleteVRF)

	h.router.HandleFunc("POST /acl", h.aclHandler.Create)
	h.router.HandleFunc("PUT /acl/{id}", h.aclHandler.Update)
	h.router.HandleFunc("DELETE /acl/{id}", h.aclHandler.Delete)
	h.router.HandleFunc("GET /acl", h.aclHandler.List)
}
