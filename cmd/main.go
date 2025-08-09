package main

import (
	"awesome-go/internal/handler"
	"awesome-go/internal/middleware"
	"awesome-go/internal/service"
	"awesome-go/pkgs/srv"
	"log"
)

func main() {
	app := srv.New()

	s := service.New()
	h := handler.New(app, s)
	m := middleware.New(s)

	h.InitializeRoutes(m)

	if err := app.Router.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
