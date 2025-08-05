package service

import (
	"github.com/pocketbase/pocketbase/core"
)

type Todo struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func (s *Service) CreateTodo(title, status string) (Todo, error) {
	todo := Todo{Title: title, Status: status}
	collection, err := s.app.FindCollectionByNameOrId("todos")
	if err != nil {
		return Todo{}, err
	}

	record := core.NewRecord(collection)
	record.Set("title", title)
	record.Set("status", status)
	err = s.app.Save(record)

	return todo, err
}

func (s *Service) DeleteTodo(id string) error {
	record, err := s.app.FindRecordById("todos", id)
	if err != nil {
		return err
	}
	return s.app.Delete(record)
}

func (s *Service) ListTodos() []Todo {
	var todos []Todo
	records, _ := s.app.FindAllRecords("todos")
	for _, r := range records {
		todos = append(todos, Todo{
			ID:     r.Id,
			Title:  r.GetString("title"),
			Status: r.GetString("status"),
		})
	}

	return todos
}
