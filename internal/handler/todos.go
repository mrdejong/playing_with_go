package handler

import (
	"awesome-go/internal/models"
	"awesome-go/views"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initializeTodos(router fiber.Router) {
	router.Get("", h.index)
	router.Post("", h.create)
	router.Delete(":id", h.delete)
}

func (h *Handler) index(c *fiber.Ctx) error {
	todos := h.service.ListTodos()
	return h.render(c, 200, views.Index(todos))
}

func (h *Handler) create(c *fiber.Ctx) error {
	var form types.TodoForm
	invalid := h.app.ParseFields(&form, c.FormValue)
	if invalid {
		return h.renderOOB(c, 200, "#todo-form", views.TodoForm(form))
	}

	todo, err := h.service.CreateTodo(form)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return h.render(c, 200, views.TodoItem(todo))
}

func (h *Handler) delete(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	h.service.DeleteTodo(idInt)
	return c.SendString("OK")
}
