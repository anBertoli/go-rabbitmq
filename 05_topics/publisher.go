package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
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

	// Start publishing messages to the exchange using different severities from all the
	// facilities. The routing key is in the format: '<facility>.<severity>'.
	wait := sync.WaitGroup{}

	wait.Add(1)
	go func() {
		generateLogs(channel, NGINX)
		wait.Done()
	}()

	wait.Add(1)
	go func() {
		generateLogs(channel, CRON)
		wait.Done()
	}()

	wait.Add(1)
	go func() {
		generateLogs(channel, SSHD)
		wait.Done()
	}()

	wait.Wait()
}

func generateLogs(channel *amqp.Channel, facility string) {
	for i := 0; ; i++ {
		routingKey := fmt.Sprintf("%s.%s", facility, randSev())
		message := fmt.Sprintf("[%s] #%d log some stuff", routingKey, i)

		err := channel.Publish(logsTopicExchange, routingKey, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
		if err != nil {
			log.Fatalf("%s", err)
		}

		log.Printf("#%d - %s", i, routingKey)
		waitRand(10)
	}
}

func waitRand(n int) {
	time.Sleep(time.Duration(rand.Intn(n)) * time.Second)
}
