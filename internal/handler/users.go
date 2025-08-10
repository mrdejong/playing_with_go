package handler

import (
	"awesome-go/internal/types"
	"awesome-go/views"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) initializeUsers(router fiber.Router) {
	router.Get("register", h.getRegister)
	router.Post("register", h.postRegister)

	router.Get("login", h.getLogin)
	router.Post("login", h.postLogin)
	router.Delete("logout", h.deleteLogin)
}

func (h *Handler) getRegister(c *fiber.Ctx) error {
	return h.render(c, 200, views.RegisterUser(types.UserForm{}))
}

func (h *Handler) postRegister(c *fiber.Ctx) error {
	var user types.UserForm
	if invalid := h.app.ParseFields(&user, c.FormValue); invalid {
		return h.render(c, http.StatusOK, views.RegistrationForm(user))
	}
	_, err := h.service.CreateUser(user)
	if err != nil {
		return h.render(c, http.StatusOK, views.RegistrationForm(user))
	}
	return h.redirect(c, "/login")
}

func (h *Handler) getLogin(c *fiber.Ctx) error {
	return h.render(c, 200, views.LoginUser(types.AuthForm{}, ""))
}

func (h *Handler) postLogin(c *fiber.Ctx) error {
	var user types.AuthForm
	if invalid := h.app.ParseFields(&user, c.FormValue); invalid {
		return h.render(c, http.StatusUnprocessableEntity, views.LoginUser(user, "Please correct the form"))
	}

	mUser, err := h.service.AuthenticateUser(user)
	if err != nil {
		return h.render(c, http.StatusUnprocessableEntity, views.LoginUser(user, "Invalid email/password provided, please correct it."))
	}

	session, err := h.service.CreateSession(mUser, string(c.Request().Header.UserAgent()), c.IP(), time.Now().Add(time.Hour*24*7))
	if err != nil {
		return h.render(c, http.StatusUnprocessableEntity, views.LoginUser(user, "Invalid email/password provided, please correct it."))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   session.ExpiresOn.Unix(),
		"token": session.Token,
	})
	tokenString, err := token.SignedString([]byte("4E4bSGLxqMYcppapKkRQrhkaNga5pcnDfSy3QDbb"))
	if err != nil {
		return err
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "auth"
	cookie.Value = tokenString
	cookie.Expires = session.ExpiresOn
	cookie.SameSite = "Lax"
	cookie.HTTPOnly = true
	cookie.Secure = true

	c.Cookie(cookie)

	return h.redirect(c, "/")
}

func (h *Handler) deleteLogin(c *fiber.Ctx) error {
	if h.authenticated(c) {
		user := h.currentUser(c)
		ip := c.IP()
		agent := string(c.Request().Header.UserAgent())

		session, err := h.service.GetSessionByMachineData(user.ID, ip, agent)
		fmt.Printf("sess: %v & err: %v", session, err)
		if err == nil {
			h.service.DeleteSession(session.ID)
		}
	}
	return h.redirect(c, "/login")
}
