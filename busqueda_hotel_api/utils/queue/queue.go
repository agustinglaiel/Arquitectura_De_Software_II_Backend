package queue

import (
	"busqueda_hotel_api/dtos"
	"busqueda_hotel_api/utils/db"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const rabbitMQURL = "amqp://user:password@localhost:5672"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func InitConsumer() {
	conn, err := amqp.Dial(rabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hotel_updates", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
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

	go func() {
		for d := range msgs {
			var hotelDTO dtos.HotelDTO
			err := json.Unmarshal(d.Body, &hotelDTO)
			if err != nil {
				fmt.Println("Error unmarshalling message:", err)
				continue
			}

			// Update Solr with the new hotel information
			doc := map[string]interface{}{
				"id":             hotelDTO.ID,
				"name":           hotelDTO.Name,
				"description":    hotelDTO.Description,
				"city":           hotelDTO.City,
				"photos":         hotelDTO.Photos,
				"amenities":      hotelDTO.Amenities,
				"room_count":     hotelDTO.RoomCount,
				"available_rooms": hotelDTO.AvailableRooms,
			}

			_, err = db.SolrClient.Update(doc, false)
			if err != nil {
				fmt.Println("Error updating Solr:", err)
				continue
			}
			fmt.Println("Successfully updated Solr with hotel:", hotelDTO.ID)
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}
