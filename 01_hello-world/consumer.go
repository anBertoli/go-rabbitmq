package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func consumer() {

	// Setting up is the same as the publisher; we open a connection and a
	// channel, and declare the queue from which we're going to consume.
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

	// Note that we declare the queue here, as well as the producer code. Because
	// we might start the consumer before the publisher, we want to make sure the
	// queue exists before we try to consume messages from it. The arguments are,
	// respectively: queue name, durable, delete when unused, exclusive, no-wait,
	// arguments.
	_, err = channel.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// The Consume method will push us messages asynchronously, so it returns a channel
	// that we can read from. The arguments are, respectively: queue name, consumer name,
	// auto ack, exclusive, no-local, no-wait, arguments.
	messages, err := channel.Consume("hello", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	for message := range messages {
		log.Printf("Received a message (routing key: '%s', exchange: '%s'): %s",
			message.RoutingKey,
			message.Exchange,
			message.Body,
		)
	}

}
