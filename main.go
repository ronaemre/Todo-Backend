package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	repository := NewTestRepository() //databse bağla
	service := NewService(repository)
	api := newAPI(&service)
	app := SetupApp(&api)
	app.Listen(":3000")
}

func SetupApp(api *API) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Get("/todos", api.GetTodosHandler)
	app.Post("/todo", api.AddTodoHandler)
	app.Delete("/todo/:id", api.DeleteTodoHandler)

	return app
}
