package main

import (
	"errors"
	"strings"

	"example.com/greetings/models"
	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

var TodoNotFound error = errors.New("Todo is Not Found")

var AlreadyExistID error = errors.New("Todo ID Already Exist")

func NewService(repository *Repository) Service {
	return Service{
		repository: repository,
	}
}

func (service *Service) GetTodo(todoId string) (*models.Todo, error) {
	todo, err := service.repository.GetTodo(todoId)

	if err != nil {
		return nil, TodoNotFound
	}

	return todo, nil
}

func (service *Service) GetTodos() ([]models.Todo, error) {
	todos, err := service.repository.GetTodos()

	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (service *Service) AddTodo(todo models.Todo) error {
	todo.ID = GenerateUUID(11)
	err := service.repository.CreateTodo(todo)

	if err != nil {
		return err
	}
	return nil

}

func (service *Service) DeleteTodo(todoId string) error {
	err := service.repository.DeleteTodo(todoId)

	if err != nil {
		return err
	}

	return nil
}

func GenerateUUID(length int) string {
	uuid := uuid.New().String()

	uuid = strings.ReplaceAll(uuid, "-", "")

	if length < 1 {
		return uuid
	}
	if length > len(uuid) {
		length = len(uuid)
	}

	return uuid[0:length]
}
