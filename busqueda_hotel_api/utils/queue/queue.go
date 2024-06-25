package queue

import (
	"busqueda_hotel_api/controllers"
	"log"

	"github.com/streadway/amqp"
)

var (
    conn *amqp.Connection
    ch   *amqp.Channel
    q    amqp.Queue
)

const rabbitMQURL = "amqp://guest:guest@localhost:5672/"

func failOnError(err error, msg string) {
    if err != nil {
        log.Panicf("%s: %s", msg, err)
    }
}

func StartReceiving() {
    var err error
    conn, err = amqp.Dial(rabbitMQURL)
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err = conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err = ch.QueueDeclare(
        "ficha_hotel-api", // name
        true,              // durable
        false,             // delete when unused
        false,             // exclusive
        false,             // no-wait
        nil,               // arguments
    )
    failOnError(err, "Failed to declare a queue")

    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    failOnError(err, "Failed to register a consumer")

    var forever chan struct{}

    go func() {
        for d := range msgs {
            log.Printf("Received a message: %s", d.Body)
            controllers.GetOrInsertByID(string(d.Body))
        }
    }()

    log.Printf("Subscription to the queue succeeded")
    <-forever
}
