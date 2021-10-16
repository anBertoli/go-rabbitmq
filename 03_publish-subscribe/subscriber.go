package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func subscriber() {

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
	// We need to declare the exchange in both publisher and subscriber, since both
	// binding a queue or publishing to a non-existing exchange will generate an error.
	err = channel.ExchangeDeclare(logsExchange, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// We want to hear about all log messages, not just a subset of them. We're also
	// interested only in currently flowing messages not in the old ones. To solve
	// that we need two things. Firstly, whenever we connect to Rabbit we need a fresh,
	// empty queue. To do this we could create a queue with a random name, or, even
	// better - let the server choose a random queue name for us. Secondly, once we
	// disconnect the consumer the queue should be automatically deleted. To do both
	// we can declare a queue with a random name (server generated), exclusive and
	// non-durable.
	//
	// As a result, when the consumer connects, a new random queue is generated. When
	// the consumer disconnects, the queue will be dropped since it's non-durable and
	// there are no other consumer (exclusive). Basically, queues will be generated
	// and destroyed dynamically when consumers connect and disconnect (and those queues
	// will be filled with flowing messages only, check also the bind below).
	queue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Now we need to tell the exchange to send messages to our queue. The relationship
	// between exchange and a queue is called a binding. From now on the logs exchange
	// will append messages to our queue.
	err = channel.QueueBind(queue.Name, "", logsExchange, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Printf("Newly generated queue '%s' binded to exchange 'logs'", queue.Name)

	// Start consuming from the new queue.
	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	for message := range messages {
		log.Printf("Received new log : [x] %s", message.Body)
	}

}
