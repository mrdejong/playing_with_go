// Package handler
package handler

import (
	"awesome-go/internal/middleware"
	"awesome-go/internal/models"
	serv "awesome-go/internal/service"
	"awesome-go/internal/types"
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

func (h *Handler) InitializeRoutes(m middleware.Middlware) {
	h.app.Router.Static("/static", "./public")

	main := h.app.Router.Group("/", m.LoadAuth)
	h.initializeHome(main)
	h.initializeUsers(main)

	todos := h.app.Router.Group("todos", m.LoadAuth, m.RequireAuth)
	h.initializeTodos(todos)
}

func (h *Handler) render(c *fiber.Ctx, status int, template templ.Component) error {
	c.Response().Header.SetStatusCode(status)
	c.Response().Header.Set("Content-Type", "text/html; charset=utf-8")

	ctx := context.WithValue(context.Background(), types.UserKey, h.currentUser(c))

	return template.Render(ctx, c.Response().BodyWriter())
}

func (h *Handler) renderOOB(c *fiber.Ctx, status int, target string, template templ.Component) error {
	c.Response().Header.SetStatusCode(status)
	c.Response().Header.Set("Content-Type", "text/html; charset=utf-8")
	c.Response().Header.Set("HX-Reswap", "outerHTML")
	c.Response().Header.Set("HX-Retarget", target)

	ctx := context.WithValue(context.Background(), types.UserKey, h.currentUser(c))
	return template.Render(ctx, c.Response().BodyWriter())
}

func (h *Handler) redirect(c *fiber.Ctx, url string) error {
	if len(c.Request().Header.Peek("HX-Request")) > 0 {
		c.Response().Header.Add("HX-Redirect", url)
		c.Response().Header.SetStatusCode(http.StatusSeeOther)
		return nil
	}
	return c.Redirect(url)
}

func (h *Handler) closeDrawer(c *fiber.Ctx) {
	c.Response().Header.Set("HX-Trigger-After-Settle", "close")
}

func (h *Handler) currentUser(c *fiber.Ctx) models.User {
	return c.UserContext().Value(types.UserKey).(models.User)
}

func (h *Handler) authenticated(c *fiber.Ctx) bool {
	return h.currentUser(c).ID > 0
}
