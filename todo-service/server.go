package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/suhail34/goGraphql-Todo/graph"
)

const defaultPort = "8080"
var rabbitMQConnection *amqp.Connection

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
  srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

  conn, err := amqp.Dial("amqp://guest:guest@rabbitmq-service.default.svc.cluster.local:5672/")
  if err != nil {
   log.Fatalf("Failed to connect to rabbitmq %v", err)
  }
  log.Println("Connection to rabbitMQ successfull")
  rabbitMQConnection = conn
  defer rabbitMQConnection.Close()
         
  http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", injectRabbitMQConnection(srv))

  log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Println(http.ListenAndServe(":"+port, nil))  
}

func injectRabbitMQConnection(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    ctx := context.WithValue(r.Context(), "rabbitMQConnection", rabbitMQConnection)
    r = r.WithContext(ctx)

    next.ServeHTTP(w, r)
  })
}
