package main

import (
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "task_queue"

func producer(taskS string) {

	// Establish the connection and create the communication channel.
	// Then we must declare a queue to send messages.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s", err)
	}

	// We need to make sure that the queue will survive a RabbitMQ node
	// restart. In order to do so, we need to declare it as durable. This
	// durable option change needs to be applied to both the producer and
	// consumer code. If you use docker queues may be deleted anyway at
	// container restart (the image must be configured properly).
	queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Inspect some basics infos about the queue and the
	// workers, before sending new messages (aka tasks).
	log.Printf("Tasks in the queue: %d\n", queue.Messages)
	log.Printf("Workers connected: %d\n", queue.Consumers)

	// Publish some messages/tasks in the queue. If we use durable queues
	// we must use the "Persistent Delivery Mode" to maintain messages in
	// the queue (transient mode will drop the messages even if the queue
	// is declared as durable).
	for i := 0; i < 30; i++ {
		taskBytes, err := json.Marshal(task{
			Name:  "hello",
			Level: i,
		})
		if err != nil {
			log.Fatalf("%s", err)
		}
		err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         taskBytes,
		})
		if err != nil {
			log.Fatalf("%s", err)
		}
		log.Printf("Sent '%s'", string(taskBytes))
		time.Sleep(time.Second)
	}

}
