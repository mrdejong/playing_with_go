package main

import (
	"awesome-go/internal/handler"
	"awesome-go/internal/service"
	"awesome-go/pkgs/srv"
	"fmt"
	"log"
	"reflect"
)

func main() {
	v := reflect.ValueOf("GTE")
	fmt.Printf("%v", v.Type())

	app := srv.New()

	s := service.New()
	h := handler.New(app, s)

	h.InitializeRoutes()

	if err := app.Router.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
