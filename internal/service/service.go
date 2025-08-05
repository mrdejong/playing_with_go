// Package service
package service

import (
	"context"

	"github.com/pocketbase/pocketbase"
)

type Service struct {
	app *pocketbase.PocketBase
}

func New(app *pocketbase.PocketBase) *Service {
	service := &Service{
		app: app,
	}
	return service
}

func (s *Service) context() context.Context {
	return context.Background()
}
