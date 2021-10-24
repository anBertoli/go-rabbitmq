package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const logsExchange = "logs"

func publisher() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Create a named exchange, called logs, of type fanout. Exchanges on one side
	// receive messages from producers and the other side they push them to queues.
	// Fanout exchanges will public messages to all bound queues. We need to declare
	// the exchange in both publisher and subscriber, since both binding a queue or
	// publishing to a non-existing exchange will generate an error.
	err = channel.ExchangeDeclare(logsExchange, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Publish to that exchange without a routing key (ignored by fanout exchanges).
	// The messages will be lost if no queue is bound to the exchange yet, but that's
	// okay for us; if no consumer is listening, yet we can safely discard the message.
	for i := 0; i < 100; i++ {
		err = channel.Publish(logsExchange, "", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("log #%d", i)),
		})
		if err != nil {
			log.Fatalf("%s", err)
		}
		time.Sleep(time.Second)
	}

}
