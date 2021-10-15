package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func producer() {

	// The connection abstracts the socket connection, and takes care of
	// protocol version negotiation and authentication and so on for us.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer conn.Close()

	// Create the communication channel.
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s", err)
	}

	// To send, we must declare a queue for us to send to; then we can
	// publish a message to the queue. The arguments are, respectively:
	// queue name, durable, delete when unused, exclusive, no-wait, arguments.
	queue, err := channel.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Inspect some basics infos about the queue,
	// before sending the message.
	log.Printf("Queue name: '%s'\n", queue.Name)
	log.Printf("Queue messages: %d\n", queue.Messages)
	log.Printf("Queue consumers: %d\n", queue.Consumers)

	// Publish some messages in the queue. The arguments are, respectively:
	// exchange, routing key, mandatory, immediate, message.
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Hello world %d", i+1)

		err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

}
