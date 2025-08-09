// Package service
package service

import (
	"awesome-go/internal/models"
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func New() *Service {
	db, err := gorm.Open(sqlite.Open("./app.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.Todo{})

	return &Service{
		db: db,
	}
}

func (s *Service) context() context.Context {
	return context.Background()
}
