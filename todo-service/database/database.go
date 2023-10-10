package database

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
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
	var userNamebase64 string = os.Getenv("USER")
  var passbase64 string = os.Getenv("PASS")
  userNameByte, _ := base64.StdEncoding.DecodeString(userNamebase64)
  passByte, _ := base64.StdEncoding.DecodeString(passbase64)
  username := string(userNameByte)
  pass := string(passByte)
  connectionString := fmt.Sprintf("mongodb+srv://%v:%v@notes.hpcfkkd.mongodb.net/?retryWrites=true&w=majority", username, pass)
	clientOptions := options.Client().ApplyURI(connectionString)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err, " can't connect to mongodb altas")
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err, " error pinging mongodb Atlas")
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

func (db *DB) CreateTodo(ctx context.Context, userId string, input *model.CreateTodoInput) (*model.Todo, error) {
	collection := db.client.Database("MyTodoService").Collection("todos")
	data := &model.Todo{
		Text:      input.Text,
		Completed: false,
		UserID:    userId,
	}
  rabbitmqContext := ctx.Value("rabbitMQConnection").(*amqp.Connection)
  ch, err := rabbitmqContext.Channel()

  if err != nil {
    log.Fatalf("Failed creating channel %v", err)
  }
  defer ch.Close()
  
  _, err = ch.QueueDeclare(
    "todo",
    true,
    false,
    false,
    false,
    nil,
  )

  if err != nil {
    log.Fatalf("Failed declaring queue %v", err)
  }

  jsonData, err := json.Marshal(data)
  if err!=nil {
    log.Fatalf("Failed to serialize JSON Data : %v", err)
  }

  err = ch.PublishWithContext(
    context.Background(),
    "",
    "todo",
    false,
    false,
    amqp.Publishing{
      ContentType: "application/json",
      Body: jsonData,
    },
  )

  if err != nil {
    log.Fatalf("Failed to publish the message %v", err)
  }
	_, err = collection.InsertOne(context.Background(), data)
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
