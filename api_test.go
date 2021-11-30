package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"example.com/greetings/models"
	"github.com/gofiber/fiber/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTodos(t *testing.T) {
	Convey("Get Todos", t, func() {
		repository := GetCleanTestRepository()
		service := NewService(repository)
		api := newAPI(&service)

		todo1 := models.Todo{
			ID:      GenerateUUID(8),
			Name:    "Deneme Todo Name",
			Details: "Deneme Todo Details ",
		}

		todo2 := models.Todo{
			ID:      GenerateUUID(8),
			Name:    "Deneme Todo Name1",
			Details: "Deneme Todo Details1 ",
		}

		todo3 := models.Todo{
			ID:      GenerateUUID(8),
			Name:    "Deneme Todo Name2",
			Details: "Deneme Todo Details2 ",
		}

		todo4 := models.Todo{
			ID:      GenerateUUID(8),
			Name:    "Deneme Todo Name3",
			Details: "Deneme Todo Details3 ",
		}

		repository.CreateTodo(todo1)
		repository.CreateTodo(todo2)
		repository.CreateTodo(todo3)
		repository.CreateTodo(todo4)

		Convey("When the get request sent", func() {
			app := SetupApp(&api)
			req, _ := http.NewRequest(http.MethodGet, "/todos", nil)

			resp, err := app.Test(req)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(resp.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then all todos should return", func() {
				actualResult := []models.Todo{}
				actualResponseBody, _ := ioutil.ReadAll(resp.Body)
				err := json.Unmarshal(actualResponseBody, &actualResult)
				So(err, ShouldBeNil)
				So(actualResult[0].ID, ShouldEqual, todo1.ID)
				So(actualResult[0].Name, ShouldEqual, todo1.Name)
				So(actualResult[0].Details, ShouldEqual, todo1.Details)
				So(actualResult[1].ID, ShouldEqual, todo2.ID)
				So(actualResult[1].Name, ShouldEqual, todo2.Name)
				So(actualResult[1].Details, ShouldEqual, todo2.Details)
				So(actualResult[2].ID, ShouldEqual, todo3.ID)
				So(actualResult[2].Name, ShouldEqual, todo3.Name)
				So(actualResult[2].Details, ShouldEqual, todo3.Details)
				So(actualResult[3].ID, ShouldEqual, todo4.ID)
				So(actualResult[3].Name, ShouldEqual, todo4.Name)
				So(actualResult[3].Details, ShouldEqual, todo4.Details)

			})

		})
	})

}

func TestAddTodo(t *testing.T) {
	Convey("Add Todo", t, func() {
		repository := GetCleanTestRepository()
		service := NewService(repository)
		api := newAPI(&service)

		todo1 := models.Todo{
			ID:      GenerateUUID(8),
			Name:    "Deneme Todo Name1",
			Details: "Deneme Todo Details1 ",
		}

		Convey("When the post request sent", func() {
			reqBody, err := json.Marshal(todo1)

			req, _ := http.NewRequest(http.MethodPost, "/todo", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Set("Content-Lenght", strconv.Itoa(len(reqBody)))

			app := SetupApp(&api)
			resp, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("When status code should be 201", func() {
				So(resp.StatusCode, ShouldEqual, fiber.StatusCreated)
			})

			Convey("Then Added todo should return", func() {
				actualResult := models.Todo{}
				actualRespBody, _ := ioutil.ReadAll(resp.Body)
				err := json.Unmarshal(actualRespBody, &actualResult)
				So(err, ShouldBeNil)
				So(actualResult, ShouldNotBeNil)
				So(actualResult.ID, ShouldEqual, todo1.ID)
				So(actualResult.Name, ShouldEqual, todo1.Name)
				So(actualResult.Details, ShouldEqual, todo1.Details)
			})

		})
	})
}

func TestDeleteTodo(t *testing.T) {
	Convey("Delete todo user wants", t, func() {

		repository := GetCleanTestRepository()
		service := NewService(repository)
		api := newAPI(&service)

		todo1 := models.Todo{
			ID:      GenerateUUID(8),
			Name:    "Deneme Todo Name1",
			Details: "Deneme Todo Details1 ",
		}
		repository.CreateTodo(todo1)

		Convey("When the delete request sent ", func() {
			app := SetupApp(&api)

			req, _ := http.NewRequest(http.MethodDelete, "/todo/"+todo1.ID, nil)
			resp, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 204", func() {
				So(resp.StatusCode, ShouldEqual, fiber.StatusNoContent)
			})

			Convey("Then todo should be deleted", func() {
				todo, err := repository.GetTodo(todo1.ID)
				So(err, ShouldNotBeNil)
				So(todo, ShouldBeNil)

			})

		})

	})

}
