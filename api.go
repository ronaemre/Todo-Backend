package main

import (
	"example.com/greetings/models"
	"github.com/gofiber/fiber/v2"
)

type API struct {
	service *Service
}

func newAPI(service *Service) API {
	return API{
		service: service,
	}
}

func (api *API) GetTodosHandler(c *fiber.Ctx) error {

	todos, err := api.service.GetTodos()

	switch err {
	case nil:
		c.JSON(todos)
		c.Status(fiber.StatusOK)

	default:
		c.Status(fiber.StatusInternalServerError)
	}

	return nil
}

func (api *API) AddTodoHandler(c *fiber.Ctx) error {

	createdTodo := models.Todo{}
	err := c.BodyParser(&createdTodo)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
	}
	err = api.service.AddTodo(createdTodo)

	switch err {
	case nil:
		c.JSON(createdTodo)
		c.Status(fiber.StatusCreated)
	}

	return nil
}

func (api *API) DeleteTodoHandler(c *fiber.Ctx) error {

	todoId := c.Params("id")
	err := api.service.DeleteTodo(todoId)

	switch err {
	case nil:
		c.Status(fiber.StatusNoContent)
	default:
		c.Status(fiber.StatusInternalServerError)
	}

	return nil
}
