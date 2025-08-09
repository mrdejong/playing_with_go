package handler

import (
	"awesome-go/internal/models"
	"awesome-go/internal/types"
	"awesome-go/views"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type NewTodo struct {
	Title  string            `form:"title"`
	Status models.TodoStatus `form:"status"`
}

func (h *Handler) initializeTodos(router fiber.Router) {
	router.Get("", h.index)
	router.Post("", h.create)
	router.Delete(":id", h.delete)
}

func (h *Handler) index(c *fiber.Ctx) error {
	todos := h.service.ListTodos()
	user := c.UserContext().Value(types.UserKey).(models.User)
	fmt.Printf("User: %v", user)
	return h.render(c, 200, views.Index(todos))
}

func (h *Handler) create(c *fiber.Ctx) error {
	var newTodo NewTodo
	err := c.BodyParser(&newTodo)
	if err != nil {
		log.Fatal(err)
		return err
	}

	todo, err := h.service.CreateTodo(newTodo.Title, newTodo.Status)
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
	return c.SendString("ok")
}
