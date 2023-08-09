package database

import (
	"context"
	"fmt"
	"log"

	"github.com/suhail34/goGraphql-Todo/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connectionString string = "mongodb+srv://skna:admin@notes.hpcfkkd.mongodb.net/?retryWrites=true&w=majority"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
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
		Todos:    []*model.Todo{},
	}
	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Fatal("Inserting Failed")
		return nil, err
	}

	return data, nil
}

func (db *DB) CreateTodo(userId, text string) (*model.Todo, error) {
	collection := db.client.Database("MyTodoService").Collection("user")
	data := &model.Todo{
		Text:      text,
		Completed: false,
		UserID:    userId,
	}
	update := bson.M{
		"$push": bson.M{
			"todos": bson.A{data},
		},
	}
	_, err := collection.UpdateOne(context.Background(), bson.M{"id": userId}, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Fatal("Todo wasn't Create")
		return nil, fmt.Errorf("Todo Wasn't Created %v", err)
	}

	return data, nil
}
