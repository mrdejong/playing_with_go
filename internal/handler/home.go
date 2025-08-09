package handler

import (
	"awesome-go/views"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initializeHome(router fiber.Router) {
	router.Get("", h.getHome)
}

func (h *Handler) getHome(c *fiber.Ctx) error {
	return h.render(c, 200, views.HomeView())
}
