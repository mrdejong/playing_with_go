package handler

import (
	"awesome-go/views"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func (h *Handler) initializeTodos(router *router.RouterGroup[*core.RequestEvent]) {
	router.GET("", h.index)
	router.POST("", h.create)
	router.DELETE("/:id", h.delete)
}

func (h *Handler) index(e *core.RequestEvent) error {
	todos := h.service.ListTodos()
	return h.render(e, 200, views.Index(todos))
}

func (h *Handler) create(e *core.RequestEvent) error {
	title := e.Request.PostFormValue("title")
	status := e.Request.PostFormValue("status")
	todo, _ := h.service.CreateTodo(title, status)
	return h.render(e, 200, views.TodoItem(todo))
}

func (h *Handler) delete(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")
	h.service.DeleteTodo(id)
	return e.String(http.StatusOK, "ok")
}
