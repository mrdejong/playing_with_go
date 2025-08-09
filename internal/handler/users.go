package handler

import (
	"awesome-go/internal/types"
	"awesome-go/views"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initializeUsers(router fiber.Router) {
	router.Get("register", h.getRegister)
	router.Post("register", h.postRegister)

	router.Get("login", h.getLogin)
	router.Post("login", h.postLogin)
}

func (h *Handler) getRegister(c *fiber.Ctx) error {
	return h.render(c, 200, views.RegisterUser(types.UserForm{}))
}

func (h *Handler) postRegister(c *fiber.Ctx) error {
	var user types.UserForm
	if invalid := h.app.ParseFields(&user, c.FormValue); invalid {
		return h.render(c, http.StatusUnprocessableEntity, views.RegisterUser(user))
	}
	_, err := h.service.CreateUser(user)
	if err != nil {
		return h.render(c, http.StatusUnprocessableEntity, views.RegisterUser(user))
	}
	return h.redirect(c, 201, "/login")
}

func (h *Handler) getLogin(c *fiber.Ctx) error {
	return h.render(c, 200, views.LoginUser(types.AuthForm{}, ""))
}

func (h *Handler) postLogin(c *fiber.Ctx) error {
	var user types.AuthForm
	if invalid := h.app.ParseFields(&user, c.FormValue); invalid {
		return h.render(c, http.StatusUnprocessableEntity, views.LoginUser(user, "Please correct the form"))
	}

	_, err := h.service.AuthenticateUser(user)
	if err != nil {
		return h.render(c, http.StatusUnprocessableEntity, views.LoginUser(user, "Invalid email/password provided, please correct it."))
	}

	return h.redirect(c, 201, "/")
}
