package service

import (
	"awesome-go/internal/models"

	"gorm.io/gorm"
)

const (
	Open       models.TodoStatus = "open"
	Pending    models.TodoStatus = "pending"
	InProgress models.TodoStatus = "in_progress"
	Completed  models.TodoStatus = "completed"
	Closed     models.TodoStatus = "closed"
)

func (s *Service) CreateTodo(title string, status models.TodoStatus) (models.Todo, error) {
	todo := models.Todo{Title: title, Status: status}
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
