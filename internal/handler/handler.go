// Package handler
package handler

import (
	serv "awesome-go/internal/service"
	"awesome-go/pkgs/srv"
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	app     *srv.Server
	service *serv.Service
}

func New(app *srv.Server, s *serv.Service) *Handler {
	return &Handler{
		app:     app,
		service: s,
	}
}

func (h *Handler) InitializeRoutes() {
	main := h.app.Router.Group("/")
	h.initializeTodos(main)
	h.initializeUsers(main)
}

func (h *Handler) render(c *fiber.Ctx, status int, template templ.Component) error {
	c.Response().Header.SetStatusCode(status)
	c.Response().Header.Set("Content-Type", "text/html; charset=utf-8")
	return template.Render(context.Background(), c.Response().BodyWriter())
}

func (h *Handler) redirect(c *fiber.Ctx, status int, url string) error {
	if len(c.Request().Header.Peek("HX-Request")) > 0 {
		c.Response().Header.Add("HX-Redirect", url)
		c.Response().Header.SetStatusCode(http.StatusSeeOther)
		return nil
	}
	c.Redirect(url, status)
	return nil
}
