package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func subscriber(severities []string) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer channel.Close()

	// Create a named exchange, called logs-routing, of type direct. Exchanges on one side
	// receive messages from producers and the other side they push them to queues. The
	// routing algorithm behind a direct exchange is simple - a message goes to the queues
	// whose binding key exactly matches the routing key of the message.
	err = channel.ExchangeDeclare(logsRoutingExchange, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// When the consumer connects, a new random queue is generated. When the consumer
	// disconnects, the queue will be dropped since it's non-durable and there are no
	// other consumer (exclusive). Basically, queues will be generated and destroyed
	// dynamically when consumers connect and disconnect.
	queue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// It is perfectly legal to bind multiple queues with the same binding key. In
	// this case messages with that routing key will be delivered to both queues.
	// It is also perfectly legal to bind a queue multiple times using different
	// binding keys. In this case, messages with one of those binding keys will
	// be delivered to the queue. Here we will bind the queue for this consumer
	// to all the provided levels of logs severity.
	for _, severity := range severities {
		err = channel.QueueBind(queue.Name, severity, logsRoutingExchange, false, nil)
		if err != nil {
			log.Fatalf("%s", err)
		}
		log.Printf(
			"queue '%s' binded to exchange '%s' with key: '%s'",
			queue.Name, logsRoutingExchange, severity,
		)
	}

	// Start consuming from the new queue. Delivered messages will be of
	// one of the bound routing keys, that is, one of the severities input.
	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}
	for message := range messages {
		log.Printf("Received new log: %s", message.Body)
	}

}
