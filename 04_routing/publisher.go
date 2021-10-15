package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

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
	defer channel.Close()

	// Create a named exchange, called logs-routing, of type direct. Exchanges on one side
	// receive messages from producers and the other side they push them to queues. The
	// routing algorithm behind a direct exchange is simple - a message goes to the queues
	// whose binding key exactly matches the routing key of the message.
	err = channel.ExchangeDeclare(logsRoutingExchange, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Start publishing messages to the exchange using different severities. The RabbitMQ
	// server will route all messages with a certain routing key to all queues bound
	// to this exchange with that binding/routing key.
	for i := 0; ; i++ {
		severity := SEVERITIES[rand.Intn(3)]
		message := fmt.Sprintf("[%s] #%d log some stuff", strings.ToUpper(severity), i)

		err = channel.Publish(logsRoutingExchange, severity, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
		if err != nil {
			log.Fatalf("%s", err)
		}

		log.Printf("#%d - %s", i, severity)
		waitRand()
	}
}

func waitRand() {
	time.Sleep(time.Duration(rand.Intn(7)) * time.Second)
}
