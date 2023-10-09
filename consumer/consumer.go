package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
  conn, err := amqp.Dial("amqp://guest:guest@rabbitmq-service.default.svc.cluster.local:5672/")
  if err != nil {
    log.Fatalf("Failed to make connection %v ", err)
  }

  defer conn.Close()

  ch, err := conn.Channel()
  if err != nil {
    log.Fatalf("Failed to create channel %v", err)
  }

  _, err = ch.QueueDeclare(
    "todo",
    true,
    false,
    false,
    false,
    nil,
  )

  if err != nil {
    log.Fatalf("Failed to declare a queue %v", err)
  }

  msgs, err := ch.Consume(
    "todo",
    "",
    true,
    false,
    false,
    false,
    nil,
  )

  if err != nil {
    log.Fatalf("Failed to register a consume %v ", err)
  }

  go func() {
    for msg := range msgs {
      log.Printf("Recieved a message : %v\n", msg.Body)
    }
  }()

  fmt.Println("Consumer is running ...")
  select {}
}
