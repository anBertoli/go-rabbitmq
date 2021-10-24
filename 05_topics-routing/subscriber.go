package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func subscriber(bindingKey string) {

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

	// Create a named exchange of type topic. Messages sent to a topic exchange can't have
	// an arbitrary routing_key - it must be a list of words, delimited by dots. The words
	// can be anything, but usually they specify some features connected to the message.
	// Example: "quick.orange.rabbit".
	//
	// The binding key must also be in the same form. The logic behind the topic exchange is
	// similar to a direct one - a message sent with a particular routing key will be delivered
	// to all the queues that are bound with a matching binding key. However, there are two
	// important special cases for binding keys:
	//
	// 		* (star): can substitute for exactly one word.
	// 		# (hash): can substitute for zero or more words.
	//
	err = channel.ExchangeDeclare(logsTopicExchange, amqp.ExchangeTopic, true, false, false, false, nil)
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

	// The routing key is in the format: '<facility>.<severity>'. Both members can
	// be a * or a #, to represent the two forms of wildcards. Based on the routing
	// key used on the producer side, this bound queue should or should not receive
	// some messages. E.g.:
	//
	// 		nginx.*		will receive all messages from the nginx facility
	//		*.error		will receive only error logs from all facilities
	//		cron.info	will receive only info logs from the cron facility
	//
	err = channel.QueueBind(queue.Name, bindingKey, logsTopicExchange, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf(
		"queue '%s' binded to exchange '%s' with binding key: '%s'",
		queue.Name, logsTopicExchange, bindingKey,
	)

	// Start consuming from the new queue. Received messages will be only the
	// ones that match the binding key used above.
	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}
	for message := range messages {
		log.Printf("Received new log: %s", message.Body)
	}

}
