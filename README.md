
# RabbitMQ concepts and usage
Examples and explanations of RabbitMQ concepts and usage in Go (taken from official RabbitMQ tutorials). Other 
related material will be probably added in the future.

# 1. Hello World

This example contains a small program that can be started in consumer o producer mode; a producer (sender) sends 
messages to a queue, and the consumer (receiver) receives messages (drains the queue). It's a "Hello World" example.

Producing means nothing more than sending. A program that sends messages is a producer. Consuming has a similar 
meaning to receiving. A consumer is a program that mostly waits to receive messages. Note that the producer, 
consumer, and broker do not have to reside on the same host; indeed in most applications they don't. An 
application can be both a producer and consumer, too.

![01 diagram](./assets/01.png)

The producer connects to the Rabbit broker (the server), declares (creates) a queue and send some messages to it.
The receiver similarly connects to the broker, declares the queue (it's an idempotent operation) and starts 
consuming the messages in the queue.

To start the example:
```shell
go run ./01_hello-world --mode producer

# in another shell
go run ./01_hello-world --mode consumer
```

#2 Worker Queues

This example creates a work queue that will be used to distribute time-consuming tasks among multiple workers.
The main idea behind work queues (aka: task queues) is to avoid doing a resource-intensive task immediately and 
having to wait for it to complete. Instead we schedule the task to be done later. We encapsulate a task as a 
message and send it to a queue. A worker process running in the background will pop the tasks and eventually 
execute the job. When you run many workers the tasks will be shared between them.

![01 diagram](./assets/02.png)

If a worker dies, we'd like the task to be delivered to another worker. In order to make sure a message is never
lost, Rabbit supports message acknowledgments. An ack(nowledgement) is sent back by the consumer to tell Rabbit
that a particular message has been received, processed and it can be deleted. If a consumer dies (or a timeout 
happen) before completing the job, no ack is sent and the task could be taken by another worker. Furthermore, 
we need to mark both the queue and messages as durable in order to avoid losing jobs if the Rabbit server dies.

Finally, another important aspect must be considered. The default fair dispatching gives one job per connected 
worker, namely it just blindly dispatches every n-th message to the n-th consumer (by default the n-th client 
will buffer all n-th messages that it can fetch). As a consequence some workers could remain more time idle than
others (think about two workers, with every odd message more intensive and time-consuming). In order to defeat 
that we can set the prefetch count with the value of 1.

A prefetch of one tells Rabbit not to give more than one message to a worker/consumer at a time. Or, in other 
words, don't dispatch a new message to a worker until it has processed and acknowledged the previous one. Instead, 
it will dispatch it to the next consumer that is not still busy (if any).

To start the example:
```shell
# start the jobs producer
go run ./02_worker-queue --mode producer

# we can start as many worker as we want, more workers means more 
# processing power and more jobs done in a period of time

# in one or more other shells 
go run ./02_worker-queue --mode worker
```

#3 Publish/Subscribe Pattern

#4 Direct Routing

#5 Topics Routing

#6 Remote Procedure Calls

#7 Publish Confirmations
