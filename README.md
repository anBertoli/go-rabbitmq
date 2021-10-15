
# RabbitMQ concepts and usage
Examples and explanations of RabbitMQ concepts and usage in Go (taken from official RabbitMQ tutorials).

# 1. Hello World

This example builds a small program that can be started in consumer o producer mode; a producer (sender) sends 
messages to a queue, and the consumer (receiver) receives messages (drains the queue). It's a "Hello World" example.

Producing means nothing more than sending. A program that sends messages is a producer. Consuming has a similar 
meaning to receiving. A consumer is a program that mostly waits to receive messages. Note that the producer, 
consumer, and broker do not have to reside on the same host; indeed in most applications they don't. An 
application can be both a producer and consumer, too.

![01 diagram](./assets/01.png)

The producer connects to the Rabbit broker (the server), declares (creates) a queue and send some messages to it.
The receiver similarly connects to the broker, declares the queue (is an idempotent operation) and starts 
consuming the messages in the queue.

```shell
go run ./01_hello-world --mode producer

# in another shell
go run ./01_hello-world --mode consumer
```

#2 Worker Queues

#3 Publish/Subscribe Pattern

#4 Direct Routing

#5 Topics Routing

#6 Remote Procedure Calls

#7 Publish Confirmations
