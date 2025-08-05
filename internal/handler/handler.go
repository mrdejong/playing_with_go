// Package handler
package handler

import (
	serv "awesome-go/internal/service"

	"github.com/a-h/templ"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

type Handler struct {
	service *serv.Service
}

func New(s *serv.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) InitializeRoutes(e *router.Router[*core.RequestEvent]) {
	main := e.Group("/todos")
	h.initializeTodos(main)
}

func (h *Handler) render(c *core.RequestEvent, status int, template templ.Component) error {
	c.Response.WriteHeader(status)
	return template.Render(c.Request.Context(), c.Response)
}
