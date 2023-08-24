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

func (db *DB) CreateUser(id, username, email string) (*model.User, error) {
	collection := db.client.Database("MyTodoService").Collection("user")
	var data = &model.User{
		ID:       id,
		Username: username,
		Email:    email,
	}
	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Fatal("Inserting Failed")
		return nil, err
	}

	return data, nil
}

func (db *DB) CreateTodo(id, userId, text string) (*model.Todo, error) {
	collection := db.client.Database("MyTodoService").Collection("todos")
	data := &model.Todo{
		ID:        id,
		Text:      text,
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
	collection := db.client.Database("MyTodoService").Collection("user")
	_ = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&data)
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
