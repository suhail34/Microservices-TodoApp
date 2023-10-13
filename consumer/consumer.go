package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
    false,
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
      var data map[string]interface{}
      err = json.Unmarshal(msg.Body, &data)
      if err!=nil {
        log.Printf("Failed to deserialize JSON Data: %v", err)
        continue
      }
      log.Printf("Recieved a message : %v\n", data)
      endTime, err := time.Parse("3:04 pm", data["endTime"].(string ))
      if err != nil {
        log.Fatalf("Error parsing endTime: %v", err)
      }
      durationRemaining := endTime.Sub(time.Now())
      if durationRemaining <= 0 {
        log.Printf("Time Duration Reached")
        if err := msg.Ack(false); err!=nil {
          log.Fatalf("Acknowledgement failed: %v", err)
        }
      }else {
        log.Printf("Time Duration Remaining")
        if err:=msg.Nack(false, true); err!=nil {
          log.Fatalf("Negative Acknowledgement failed:%v", err)
        }
      }
    }
  }()

  fmt.Println("Consumer is running ...")
  select {}
}
