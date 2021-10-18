
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

# 2. Worker Queues

This example creates a work queue that will be used to distribute time-consuming tasks among multiple workers.
The main idea behind work queues (aka: task queues) is to avoid doing a resource-intensive task immediately and 
having to wait for it to complete. Instead we schedule the task to be done later. We encapsulate a task as a 
message and send it to a queue. A worker process running in the background will pop the tasks and eventually 
execute the job. When you run many workers the tasks will be shared between them.

![02 diagram](./assets/02.png)

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
go run ./02_workers-queue --mode producer

# we can start as many worker as we want, more workers means more 
# processing power and more jobs done in a period of time

# in one or more other shells 
go run ./02_workers-queue --mode worker
```

# 3. Publish/Subscribe 

In previous parts of the tutorial we sent and received messages to and from a queue. In the full messaging model 
in Rabbit there are also exchanges. The core idea in the messaging model in RabbitMQ is that the producer never 
sends any messages directly to a queue. Actually, quite often the producer doesn't even know if a message will 
be delivered to any queue at all. Instead, the producer can only send messages to an exchange. Then, Rabbit sent
those messages to one or more queues bound to that exchange.

An exchange is a very simple thing. On one side it receives messages from producers and the other side it pushes
them to queues. The exchange must know exactly what to do with a message it receives. The rules for exchanges are 
defined by their type (direct, topic, headers and fanout).

This example is a logging system, we want to send all logs produced to an exchange and then to all 'subscribers'. 
To do this we will use an exchange of type 'fanout' which just broadcasts all the messages it receives to all 
the queues it knows. When a consumer/subscriber joins we create a new empty queue with a random name, a queue that
is specific for that subscriber. Then we bind the queue to the 'logs' exchange. Note that the subscriber is not 
interested in messages sent before it is connected to the server. Secondly, when the subscriber disconnects, 
the queue will be dropped.

Basically, queues will be generated and destroyed dynamically when consumers connect and disconnect (and those
queues will be filled with flowing messages only).

![03 diagram](./assets/03.png)

To start the example:
```shell
# start the jobs producer
go run ./03_publish-subscribe --mode publisher

# we can start as many subscribers as we want, 
# similarly to a real subscription system

# in one or more other shells 
go run ./03_publish-subscribe --mode subscriber
```

# 4. Direct Routing

In this tutorial we're going to make it possible to subscribe only to a subset of the messages. The structure is
similar to the previous program, but we use different routing/binding keys in both queue bindings and message sending. 
To build this example we use a 'direct' exchange instead. The routing algorithm behind a direct exchange is simple - 
a message goes to the queues whose binding key exactly matches the routing key of the message.

![04 diagram](./assets/04.png)

In the setup in figure, the direct exchange X has two queues bound to it. The first queue is bound with binding 
key _orange_, and the second has two bindings, one with binding key _black_ and the other one with _green_. Here,
a message published to the exchange with a routing key orange will be routed to the first queue. Messages with a 
routing key of black or green will go to second queue. All other messages will be discarded. It is perfectly legal 
to bind multiple queues with the same binding key. In this case a message with that routing key is sent to all 
queues bound with that key.

This example is a logging system, where the logs' producer send logs to an exchange with different severities. 
Subscribers/consumers can listen for logs with multiple severities binding their freshly-created and exclusive 
queues with one or more keys. While consuming their queues, they will receive the subset of messages they 
subscribed for.

# 5. Topics Routing

There is another type of exchange named 'topic'. With these type of exchange we can bind a queue with multiple
criteria. Messages sent to a topic exchange can't have an arbitrary routing key - it must be a list of words,
delimited by dots (e.g. "quick.orange.rabbit"). The matching rules are similar to those of the 'direct' exchange
with some differences: a star (*) can substitute for exactly one word, a hash (#) can substitute for multiple 
words. When special characters "*" (star) and "#" (hash) aren't used in bindings, the topic exchange will behave 
just like a direct one.

![05 diagram](./assets/05.png)

This simple program is a more advanced logging system where the both severities and the source form the routing key.
In the form _<facility>.<severity>_. In the consumer both members can be a * or a #, to represent the two forms of
wildcards. Based on the routing key used on the producer side, bound queues should or should not receive some 
messages. E.g.:

- nginx.*       will receive all messages from the nginx facility
- *.error       will receive only error logs from all facilities
- cron.info	    will receive only info logs from the cron facility


# 6. Remote Procedure Calls

![06 diagram](./assets/06.png)

# 7. Publish Confirmations

![07 diagram](./assets/07.png)
