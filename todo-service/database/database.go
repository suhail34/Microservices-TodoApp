package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/suhail34/goGraphql-Todo/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var connectionString string = os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(connectionString)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB Atlas Connected")
	return &DB{
		client: client,
	}
}

func (db *DB) CreateUser(input *model.CreateUserInput) (*model.User, error) {
	collection := db.client.Database("MyTodoService").Collection("user")
	var data = &model.User{
		Username: input.Username,
		Email:    input.Email,
	}
	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Fatal("Inserting Failed")
		return nil, err
	}

	return data, nil
}

func (db *DB) CreateTodo(userId string, input *model.CreateTodoInput) (*model.Todo, error) {
	collection := db.client.Database("MyTodoService").Collection("todos")
	data := &model.Todo{
		Text:      input.Text,
		Completed: false,
		UserID:    userId,
	}
	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Fatal("Todo wasn't Create")
		return nil, fmt.Errorf("Todo Wasn't Created %v", err)
	}

	return data, nil
}

func (db *DB) GetUser(id string) (*model.User, error) {
	var data *model.User
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("Invalid ID", err)
	}
	collection := db.client.Database("MyTodoService").Collection("user")
	_ = collection.FindOne(context.Background(), bson.M{"_id": _id}).Decode(&data)
	if data == nil {
		return nil, fmt.Errorf("User not present")
	}
	return data, nil
}

func (db *DB) GetUserTodos(userId string) ([]*model.Todo, error) {
	var todos []*model.Todo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.client.Database("MyTodoService").Collection("todos")
	cursor, err := collection.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		log.Fatal("Cannot find any entry ", err)
	}
	for cursor.Next(ctx) {
		var todo *model.Todo
		if err := cursor.Decode(&todo); err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (db *DB) GetTodo() ([]*model.Todo, error) {
	var todos []*model.Todo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.client.Database("MyTodoService").Collection("todos")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("No todo found with specified ID")
	}
	for cursor.Next(ctx) {
		var todo *model.Todo
		if err := cursor.Decode(&todo); err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (db *DB) UpdateTodo(id, userId string, input *model.UpdateTodoInput) (*model.Todo, error) {
	var todo *model.Todo
	collection := db.client.Database("MyTodoService").Collection("todos")
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("Invalid ID", err)
	}
	updateFields := bson.M{}
	if input.Text != nil {
		updateFields["text"] = input.Text
	}
	if input.Completed != nil {
		updateFields["completed"] = input.Completed
	}
	update := bson.M{
		"$set": updateFields,
	}
	filter := bson.M{"_id": _id, "userId": userId}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal("Failed updating todo", err)
	}
	err = collection.FindOne(context.Background(), filter).Decode(&todo)
	if err != nil {
		log.Fatal("Error finding todo for user", err)
	}
	return todo, nil
}

func (db *DB) DeleteTodo(id string) (*model.Todo, error) {
	var todo *model.Todo
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("Invalid ID", err)
	}
	collection := db.client.Database("MyTodoService").Collection("todos")
	err = collection.FindOne(context.Background(), bson.M{"_id": _id}).Decode(&todo)
	if err != nil {
		return nil, fmt.Errorf("No Todo Item present with that id %v", err)
	}
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": _id})
	if err != nil {
		return nil, fmt.Errorf("Delete operation failed %v", err)
	}
	return todo, nil
}
