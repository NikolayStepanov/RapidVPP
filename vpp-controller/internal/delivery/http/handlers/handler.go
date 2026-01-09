package handlers

import (
	"net/http"
)

type Handler struct {
	router *http.ServeMux
}

func NewHandler() *Handler {
	handler := &Handler{
		router: http.NewServeMux(),
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
