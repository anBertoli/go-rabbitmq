package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func start() {

	// We start by establishing the connection, the channel and declaring the queue.
	// Then, we declare a queue where we are testing publisher confirms.
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

	queue, err := channel.QueueDeclare("test", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Publisher confirms are a RabbitMQ extension to the AMQP 0.9.1 protocol, so
	// they are not enabled by default. Publisher confirms are enabled at the
	// channel level with the Confirm() method. This method must be called
	// on every channel that you expect to use publisher confirms. Confirms
	// should be enabled just once, not for every message published.
	err = channel.Confirm(false)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Let's start with the simplest approach to publishing with confirms, that is,
	// publishing a message and waiting synchronously for its confirmation
	confirmations := channel.NotifyPublish(make(chan amqp.Confirmation, 100))

	// Send some messages and wait synchronously the broker confirmation.
	// The confirmation is performed serially, that is, after each message
	// we wait the confirmation (so we don't batch published messages).
	for i := 0; i < 10; i++ {
		err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
			ContentType: "plain/text",
			Body:        []byte("abc"),
		})
		if err != nil {
			log.Fatalf("%s", err)
		}

		waitForConfirm(confirmations, 1)
	}

	// Now we will send messages in batch and we'll wait for confirmations in
	// batches. After every "batch-size" messages sent we want to receive the
	// broker confirmation for all of them. Waiting for a batch of messages to
	// be confirmed improves throughput drastically over waiting for
	// individual messages confirmations.
	msgsInFlight := 0
	batchSize := 50

	for i := 0; i < 1000; i++ {
		err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
			ContentType: "plain/text",
			Body:        []byte("abc"),
		})
		if err != nil {
			log.Fatalf("%s", err)
		}

		// If we didn't reach the batch size we didn't wait confirmations,
		// otherwise we wait for them. We expect to receive a number of
		// confirmations equal to the batch size.
		msgsInFlight++
		if msgsInFlight%batchSize == 0 {
			waitForConfirm(confirmations, msgsInFlight)
			msgsInFlight = 0
		}
	}

	// Wait for pending confirms.
	if msgsInFlight > 0 {
		waitForConfirm(confirmations, msgsInFlight)
	}

}

// Utility function that waits for n message confirmations
// from the provided confirmation channel.
func waitForConfirm(confirmations <-chan amqp.Confirmation, n int) {
	// Allow at most five seconds for
	// all the broker confirmations.
	timer := time.NewTimer(5 * time.Second)
	for i := 0; i < n; i++ {
		select {
		case cnf := <-confirmations:
			fmt.Printf("confirmation %+v\n", cnf)
			if !cnf.Ack {
				log.Fatal("not acked")
			}
		case <-timer.C:
			log.Fatal("timeout")
		}
	}
}
