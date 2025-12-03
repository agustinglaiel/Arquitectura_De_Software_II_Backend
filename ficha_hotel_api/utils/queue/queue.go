package queue

import (
	"context"
	"encoding/json"
	"ficha_hotel_api/dtos"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Init() {
	var err error
	var err2 error
	
	rabbitMQURL := "amqp://guest:guest@rabbit:5672/"
	/*rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@rabbit:5672/"
	}*/

	conn, err = amqp.Dial(rabbitMQURL)
	if err!=nil{
		log.Println("configuracion de docker no existente, Probando configuracion local")
		rabbitMQURL2 := "amqp://guest:guest@localhost:5672/"
		conn, err2 = amqp.Dial(rabbitMQURL2)
		if err2 != nil{
			log.Println("Error en local:",err2)
			failOnError(err, "Failed to connect to RabbitMQ Dokcer")
		}
	}

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err = ch.QueueDeclare(
		"ficha_hotel-api", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue")
}

func Send(id string, action string) {
	if ch == nil {
		log.Println("Channel is not initialized")
		return
	}

	// Prepare a message in the same format as the consumer expects
	queueDto := dtos.QueueDto{
		Id:     id,
		Action: action,
	}

	body, err := json.Marshal(queueDto)
	if err != nil {
		log.Fatalf("Error encoding JSON: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json", // Change content type to JSON
			Body:         body,
		})

	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	} else {
		log.Printf(" [x] Sent message: ID %s, Action %s\n", id, action)
	}
}
