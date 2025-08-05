package main

import (
	"awesome-go/internal/handler"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()
	h := handler.New(nil)

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		h.InitializeRoutes(e.Router)

		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
