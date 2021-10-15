package main

import (
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func worker() {

	// Establish the connection and create the communication channel,
	// then declare the queue.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Declare the queue with the parameter durable and start consuming
	// messages from it.
	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Before consuming messages we need to set the prefetching values. Currently,
	// the fair dispatching gives one job per connected worker, namely it just
	// blindly dispatches every n-th message to the n-th consumer. As a consequence
	// some workers could remain more time idle than others. In order to defeat that
	// we can set the prefetch count with the value of 1 (by default the n-th client
	// will buffer all n-th messages that it can fetch).
	//
	// A prefetch of one tells RabbitMQ not to give more than one message to a worker
	// at a time. Or, in other words, don't dispatch a new message to a worker until
	// it has processed and acknowledged the previous one. Instead, it will dispatch
	// it to the next consumer that is not still busy.
	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("%s", err)
	}

	messages, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// We don't want to lose any tasks. If a worker dies, we'd like the task to be
	// delivered to another worker. To do this, we use acknowledgements (ack). They
	// are messages sent back by the consumer to tell RabbitMQ that a particular
	// message has been received, processed and that RabbitMQ is free to delete it.
	//
	// If no ack is sent rabbitMQ will re-enqueue the message. If there are other
	// consumers online at the same time, it will then quickly redeliver it to
	// another consumer.
	//
	// To acknowledge manually the messages we must consume from the queue/channel
	// with the auto-ack parameter set to false, otherwise the RabbitMQ server
	// will automatically delete them after sending them.
	for message := range messages {
		log.Printf("Received a message (routing key: '%s', exchange: '%s'): %s",
			message.RoutingKey,
			message.Exchange,
			message.Body,
		)

		// Parse the JSON-formatted rabbit message into a worker
		// task and simulate working on the task for some seconds.
		task, err := parseTaskMessage(message)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Task in progress: %s\n", message.Body)
		time.Sleep(time.Duration(task.Level) * time.Second)

		// The 'multiple' argument dictate if the ack should only for
		// this message or should be a collective ack (i.e. all message
		// sent via this channel are acked).
		err = message.Ack(false)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Task completed: %s\n", message.Body)
	}
}

func parseTaskMessage(message amqp.Delivery) (t task, err error) {
	err = json.Unmarshal(message.Body, &t)
	return
}
