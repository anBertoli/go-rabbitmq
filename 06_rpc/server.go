package main

import (
	"log"
	"strconv"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

func startServers() {
	rpcServers := sync.WaitGroup{}

	// Start 3 rpc servers at the application level (3 goroutines). Wait
	// until they are closed (never, in the current implementation).
	for i := 0; i < 3; i++ {
		rpcServers.Add(1)
		go func(id int) {
			defer rpcServers.Done()
			server(id)
		}(i)
	}

	rpcServers.Wait()
}

func server(id int) {

	// We start by establishing the connection,
	// the channel and declaring the queue.
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

	// Declare the RPC work queue. We will drain from this shared queue, and we
	// will put responses on the client-dedicated response queues.
	queue, err := channel.QueueDeclare(rpcQueue, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// The server is assigned only one task at a time. This will be useful
	// if we run more instances of RPC servers (which actually are workers).
	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Drain RPC messages from the work queue. We will drain messages/tasks from
	// this shared queue, and we will put responses on the client-dedicated response
	// queue (the 'callback' queue). The correlation-id field is used to correlate
	// the response with its RPC request, while the reply-to field is used to know
	// where we must put the RPC response message.
	rpcRequests, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	for rpcRequest := range rpcRequests {
		log.Printf("[%d] Received new RPC request: %s\n", id, rpcRequest.Body)
		num, err := strconv.Atoi(string(rpcRequest.Body))
		if err != nil {
			log.Fatalf("%s", err)
		}

		// Execute the task and publish the response to the callback queue
		// in the "reply-to" field, with the proper correlation id.
		rpcResponse := []byte(strconv.Itoa(fib(num)))

		err = channel.Publish("", rpcRequest.ReplyTo, false, false, amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: rpcRequest.CorrelationId,
			Body:          rpcResponse,
		})
		if err != nil {
			log.Fatalf("%s", err)
		}

		// Finally, ack the RPC request in order to be deleted from the
		// shared RPC requests queue.
		err = rpcRequest.Ack(false)
		if err != nil {
			log.Fatalf("%s", err)
		}
	}
}

// Simulate a task with this function.
func fib(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}
