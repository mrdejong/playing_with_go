// Package models
package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

type Session struct {
	gorm.Model
	UserID    uint
	User      User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Token     string
	ExpiresOn time.Time
	UserAgent string
	IP        string
}

type TodoStatus string

type Todo struct {
	gorm.Model
	Title  string     `json:"title"`
	Status TodoStatus `json:"status"`
}
