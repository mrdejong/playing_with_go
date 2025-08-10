package service

import (
	"awesome-go/internal/models"
	"awesome-go/internal/types"

	"gorm.io/gorm"
)

const (
	Open       models.TodoStatus = "open"
	Pending    models.TodoStatus = "pending"
	InProgress models.TodoStatus = "in_progress"
	Completed  models.TodoStatus = "completed"
	Closed     models.TodoStatus = "closed"
)

// title string, status models.TodoStatus
func (s *Service) CreateTodo(form types.TodoForm) (models.Todo, error) {
	todo := models.Todo{Title: form.Title.Value, Status: models.TodoStatus(form.Status.Value)}
	err := gorm.G[models.Todo](s.db).Create(s.context(), &todo)
	return todo, err
}

func (s *Service) DeleteTodo(id int) error {
	_, err := gorm.G[models.Todo](s.db).Where("id = ?", id).Delete(s.context())
	return err
}

func (s *Service) ListTodos() []models.Todo {
	todos, _ := gorm.G[models.Todo](s.db).Find(s.context())
	return todos
}
