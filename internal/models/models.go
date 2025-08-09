package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

type TodoStatus string

type Todo struct {
	gorm.Model
	Title  string     `json:"title"`
	Status TodoStatus `json:"status"`
}
