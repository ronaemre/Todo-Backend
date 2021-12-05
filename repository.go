package main

import (
	"context"
	"log"
	"time"

	"example.com/greetings/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoEntity struct {
	ID      string `bson:"id"`
	Name    string `bson:"name"`
	Details string `bson:"details"`
}

func convertTodoEntityToModel(todoEntity TodoEntity) models.Todo { //bson-json
	return models.Todo{
		ID:      todoEntity.ID,
		Name:    todoEntity.Name,
		Details: todoEntity.Details,
	}
}

type Repository struct {
	client *mongo.Client
}

func (repository *Repository) CreateTodo(todo models.Todo) error {

	collection := repository.client.Database("TodoTestDatabase").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, todo)
	if err != nil {
		return err
	}
	return nil

}

func (repository *Repository) GetTodos() ([]models.Todo, error) {
	collection := repository.client.Database("TodoTestDatabase").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{}

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	todos := []models.Todo{}
	for cur.Next(ctx) {
		var todo models.Todo
		err := cur.Decode(&todo)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (repository *Repository) GetTodo(todoId string) (*models.Todo, error) {
	collection := repository.client.Database("TodoTestDatabase").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": todoId}

	cur := collection.FindOne(ctx, filter)
	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, TodoNotFound
	}

	todoEntity := TodoEntity{}
	err := cur.Decode(&todoEntity)

	if err != nil {
		return nil, err
	}

	todo := convertTodoEntityToModel(todoEntity)

	return &todo, nil

}

func (repository *Repository) DeleteTodo(todoId string) error {
	collection := repository.client.Database("TodoTestDatabase").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": todoId}

	_, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	return nil

}

func NewRepository() *Repository {
	uri := "mongodb+srv://rona:rona123@clustertodo.emo1b.mongodb.net/myFirstDatabase?authSource=admin&replicaSet=atlas-frb3e6-shard-0&w=majority&readPreference=primary&appname=MongoDB%20Compass&retryWrites=true&ssl=true"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return &Repository{client}
}

func NewTestRepository() *Repository { //mongo localae bağlama test için
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer cancel()
	client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return &Repository{client}
}

func GetCleanTestRepository() *Repository { //her test koştuğunda database sıfırla

	repository := NewTestRepository()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	todoDB := repository.client.Database("TodoTestDatabase")
	todoDB.Drop(ctx)

	return repository
}
