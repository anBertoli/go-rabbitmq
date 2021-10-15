package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func client() {

	// We start by establishing the connection, channel and declaring the queue.
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

	// Declare a random named queue exclusive for this client and start
	// consuming from this queue. This phase is basically a setup of the
	// client-exclusive callback queue. When the client disconnects the
	// queue is dropped (durable = false, auto-delete = true), so all
	// responses are discarded.
	callbackQueue, err := channel.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Auto ack and exclusive on. The RPC responses are automatically acknowledged
	// and the queue must have only this consumer.
	rpcResponses, err := channel.Consume(callbackQueue.Name, "", true, true, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Publish RPC messages to the work queue. The server process will drain from
	// this shared queue and puts responses on the client-dedicated response queues.
	// The correlation-id field is used to correlate the response with its RPC
	// request, while the reply-to field is used to inform the RPC server about
	// where to send the response.
	//
	// Note: this is usually inefficient, it's probably better to send RPC requests
	// in batch, recording all correlation IDs and check this list while receiving
	// responses.
	for {
		rpcRequest := strconv.Itoa(randInt(5, 15))
		rpcCorrelationId := randomString(32)

		err = channel.Publish("", rpcQueue, false, false, amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: rpcCorrelationId,
			ReplyTo:       callbackQueue.Name,
			Body:          []byte(rpcRequest),
		})
		if err != nil {
			log.Fatalf("%s", err)
		}
		log.Printf("Sent RPC request: %s\n", rpcRequest)

		// Wait for the RPC response, we expect to find a message on the client-exclusive
		// response queue. The correlation ID must match the value we used in the same field
		// while publishing the message. Although unlikely, it is possible that the RPC server
		// will die just after sending us the answer, but before sending an acknowledgment
		// message for the request. If that happens, the restarted RPC server will process the
		// request again. That's why on the client we must handle the duplicate responses
		// gracefully, and the RPC should ideally be idempotent.
		for rpcResponse := range rpcResponses {
			if rpcCorrelationId != rpcResponse.CorrelationId {
				continue
			}
			res, err := strconv.Atoi(string(rpcResponse.Body))
			if err != nil {
				log.Fatalf("%s", err)
			}
			log.Printf("RPC response: [ %d ]\n", res)
			break
		}

		time.Sleep(time.Second * 2)
	}

}

// Utility functions.
func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
