package queue

import (
	config "busqueda_hotel_api/config"
	"busqueda_hotel_api/controllers"
	"busqueda_hotel_api/dtos"
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var (
	conn *amqp.Connection
)

//const rabbitMQURL = "amqp://user:password@localhost:5672/"

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func StartReceiving() {
	QueueConn, err := amqp.Dial(config.AMPQConnectionURL)
		if err!=nil{
			QueueConn, err = amqp.Dial(config.AMPQConnectionURLlocal)
			failOnError(err, "Can't connect to AMQP")
		}

	defer QueueConn.Close()

	amqpChannel, err := QueueConn.Channel()
	failOnError(err, "Can't create a amqpChannel")
	defer amqpChannel.Close()

	err = amqpChannel.Qos(1, 0, false)
	failOnError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		config.QUEUENAME,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Could not register consumer")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)

			var queueDto dtos.QueueDto

			err := json.Unmarshal(d.Body, &queueDto)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			log.Printf("ID %s, Action %s", queueDto.Id, queueDto.Action)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

			if queueDto.Action == "INSERT" || queueDto.Action == "UPDATE" {

				if queueDto.Action == "UPDATE" {

					err := controllers.Delete(queueDto.Id)

					if err != nil {
						failOnError(err, "Error deleting from Solr")
					}

				}

				err := controllers.AddFromId(queueDto.Id)

				if err != nil {
					failOnError(err, "Error inserting or deleting from Solr")
				}

			} else if queueDto.Action == "DELETE" {
				err := controllers.Delete(queueDto.Id)

				if err != nil {
					failOnError(err, "Error inserting or deleting from Solr")
				}
			}

		}
	}()

	// Stop for program termination
	<-stopChan
}
