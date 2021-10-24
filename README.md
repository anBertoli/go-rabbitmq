
# RabbitMQ concepts and usage
Examples and explanations of RabbitMQ concepts and usage in Go. Taken from the official RabbitMQ tutorials, 
with extended examples and additional clarifications.

In all examples we assume that a RabbitMQ instance is up and running at `127.0.0.1:5672`. The `rabbit.sh` file 
contains a script to start such an instance in a Docker image. Note that all the things that could be persisted 
by RabbitMQ are cleaned by this Docker image when it stops (e.g. persistent queues).

# 1. Hello World

This example contains a small program that can be started in consumer o producer mode; a **producer** (sender) sends 
messages to a queue, and a **consumer** (receiver) receives messages from a queue. It's a "Hello World" example.

Producing means nothing more than sending. A program that sends messages is a producer. Consuming has a similar 
meaning to receiving. A consumer is a program that mostly waits to receive messages. Note that the producer, the
consumer, and the message broker do not have to reside on the same host; indeed in most applications they don't. 
An application can be both a producer and consumer too.

![01 diagram](./assets/01_new.png)

In this simple example, the producer connects to the Rabbit message broker (the server), declares (creates) a queue
and send some messages to it. The receiver similarly connects to the broker, declares the queue and starts consuming
the messages from the queue. Declaring a queue is an idempotent operation and it's a good practice to perform it on 
both the consumer and the producer side, to ensure that the queue exists upon connection. 

To start the example:
```shell
# Start the producer.
go run ./01_hello-world --mode producer

# In one or more other shells, start the consumers.
go run ./01_hello-world --mode consumer
```

# 2. Worker Queues

This example creates a work queue that will be used to distribute time-consuming tasks among multiple workers. The
main idea behind work queues (aka task queues) is to avoid performing resource-intensive tasks immediately and have
to wait for them. Instead, we want to schedule the task to be done later. We encapsulate a task as a message and 
send it to a queue. A running worker process will pop the tasks/messages and execute the job. When you run many
workers at the same time the tasks will be equally distributed among them.

![02 diagram](./assets/02_new.png)

If a worker dies, we'd like the task to be delivered to another worker. In order to make sure a message is never
lost, Rabbit supports **message acknowledgments**. An ack(nowledgement) is sent back by the consumer to tell Rabbit
that a particular message has been received and processed, so it can be deleted safely. If a consumer dies (or a 
timeout happens) before completing the job, no ack is sent and the task could be taken by another worker. Furthermore, 
we need to mark both the queue and messages as durable in order to avoid losing jobs if the Rabbit server dies.

Finally, another important aspect must be considered. The default _fair dispatching_ of messages assigns one job per 
connected worker, that is, it just blindly dispatches every _n-th message to the n-th consumer_ (and the n-th client
will buffer all n-th messages that it can fetch and store). As a consequence some workers could remain more time idle
than others. Think about two workers connected and every odd message being more intensive and time-consuming. In this
scenario the second worker is assigned the more intensive jobs, while the first one will remain generally idle for 
more time. In order keep all the workers busy in an equal way, we can set the **prefetch count** to the value of 1.

A prefetch of one tells Rabbit not to give more than one message to a worker/consumer at a time. More precisely, 
don't dispatch a new message to a worker until it has processed and acknowledged the previous one. Instead, Rabbit 
will dispatch the next message to the next consumer that is not still busy (if any).

To start the example:
```shell
# Start the jobs producer.
go run ./02_workers-queue --mode producer

# We can start as many worker as we want, more workers means more 
# processing power and more jobs done in a period of time. Run in
# one or more other shells: 
go run ./02_workers-queue --mode worker
```

# 3. Publish/Subscribe 

In previous examples we sent and received messages to and from a queue. In the full messaging model of Rabbit there
are also **exchanges**. The core idea of the RabbitMQ messaging model is that the producer never sends any messages
directly to a queue. Actually, quite often the producer doesn't even know if a message will be delivered to any
queue at all. Instead, the producer can only send messages to an exchange. Then, Rabbit (which is a message broker 
indeed) will send those messages to one or more queues _bound_ to that exchange.

An exchange is a very simple thing. On one side it receives messages from producers and on the other side it pushes
them to queues. The exchange must know exactly what to do with a message it receives. The rules for exchanges are 
defined by their type (**direct**, **topic**, **headers** and **fanout**).

![03 diagram](./assets/03_new.png)

This example is a logging system that implements the publisher/subscriber design pattern. The _publisher_ will send
all the produced logs to an exchange, then they will be routed to all _subscribers_ (the consumers). To do this we
will use an exchange of type **_fanout_** which just broadcasts all the messages it receives to all the queues it knows.

When a consumer/subscriber joins we create a new empty queue with a random name, a queue that is specific for that
subscriber. Then we bind the queue to the `logs` exchange. Note that the subscriber is not interested in messages
sent before it is connected to the server. Secondly, when the subscriber disconnects, the queue will be dropped.
Basically, queues will be generated and destroyed dynamically when consumers connect and disconnect (and those
queues will be filled with flowing messages only).

To start the example:
```shell
# Start the publisher.
go run ./03_publish-subscribe --mode publisher

# We can start as many subscribers as we want, 
# similarly to a real subscription system. Run
# in one or more other shells: 
go run ./03_publish-subscribe --mode subscriber
```

# 4. Direct Routing

In this example we're build a logging system as before, but we are going to make it possible to subscribe only to a 
subset of messages sent by a publisher. The structure is similar to the previous example, but we'll use different 
**routing/binding keys** in both queue bindings and message sending. 

Routing keys are strings used by producers to control where messages must be sent, specifically those messages are 
sent to the queues that are bound to the exchange with a compatible binding key (binding key is like the routing
key, but for consumers). Different types of exchanges have different types of matching criteria. The _fanout_ 
exchange used in the previous example just send all messages to all bound queues, regardless of the keys used. To
build this example we use a **_direct_** exchange instead. The routing algorithm behind a direct exchange is simple: 
a message goes to the queues whose binding key exactly matches the routing key of the message.

![04 diagram](./assets/04_new.png)

In the figure above, the direct exchange has three queues bound to it. The first queue is bound with the binding 
key `info`, the second queue is bound with the binding key `warn` and the third queue has two bindings, one with
the binding key `warn` and the other one with `error`. Here, a message published to the exchange with a routing key 
info will be routed to the first queue. Messages with a routing key of _warn_ will go to second and third queues.
Messages with a routing key of _error_ will go to the third queue only. All other messages will be discarded. It 
is perfectly legal to bind multiple queues with the same binding key. 

In the example the producer send logs to an exchange with different severities (that are the routing keys of the 
system). Subscribers/consumers can listen for logs with one or more severities, binding their freshly-created and
exclusive queues with the related keys. While consuming their queues, consumers will receive only the subset of 
messages they subscribed for.

To start the example:
```shell
# Start the logs producer (the publisher).
go run ./04_direct-routing --mode publisher

# Start a subscriber for info logs, in a new shell.
go run ./04_direct-routing --mode subscriber --sevs info

# Start a subscriber for warn logs, in a new shell.
go run ./04_direct-routing --mode subscriber --sevs warn

# Start a subscriber for warn and error logs, in a new shell. 
go run ./04_direct-routing --mode subscriber --sevs warn-error
```

# 5. Topics Routing

There is another type of exchange named 'topic'. With these type of exchange we can bind a queue with multiple
criteria. Messages sent to a topic exchange can't have an arbitrary routing key - it must be a list of words,
delimited by dots (e.g. "quick.orange.rabbit"). The matching rules are similar to those of the 'direct' exchange
with some differences: a star (\*) can substitute for exactly one word, a hash (#) can substitute for multiple 
words. When special characters "*" (star) and "#" (hash) aren't used in bindings, the topic exchange will behave 
just like a direct one.

![05 diagram](./assets/05_new.png)

This simple program is a more advanced logging system where the both severities and the source form the routing key.
In the form _<facility>.<severity>_. In the consumer both members can be a * or a #, to represent the two forms of
wildcards. Based on the routing key used on the producer side, bound queues should or should not receive some 
messages. E.g.:

- nginx.*       will receive all messages from the nginx facility
- *.error       will receive only error logs from all facilities
- cron.info	    will receive only info logs from the cron facility

To start the example:
```shell
# start the logs producer
go run ./05_topics --mode publisher

# to start a subscriber that listens for errors from all sources 
go run ./05_topics --mode subscriber --bind *.error

# to start a subscriber for info logs from nginx
go run ./05_topics --mode subscriber --bind nginx.info
```

# 6. Remote Procedure Calls

This example is a simple setup of a pattern called Remote Procedure Call (RPC). With this pattern we want to run
a function on a remote computer, collecting the response on the client side. The system is composed by a client 
and one or more RPC servers. In the example a client sends a request message to a shared work queue (the RPC requests
queue), the servers reads the task message, performs the job, and replies with a RPC response message in a 
consumer-exclusive callback queue.

![06 diagram](./assets/06_new.png)

In order to receive the response in a specific queue we need to send the callback queue address with the request. 
Callback queues are exclusive and generated upon the client connection. RPC responses are auto-acked from the queue 
while RPC requests must be acknowledged by the server when it finishes working on that task. There is another field, 
the correlation-id used to correlate the response with its related RPC request. The reply-to field is used instead 
to inform the RPC server about where to send the response.

Note that we send one RPC request per time and we wait for the response sequentially. It's probably better to send
RPC requests in batch, recording all correlation IDs and check this list while receiving responses (not done here, 
since we are interested in the architecture not in performance).

To start the example:
```shell
# start the RPC server (actually, three servers are started)
go run ./06_rpc --mode server

# start one or more clients in different shells
go run ./06_rpc --mode client
```

# 7. Publish Confirmations

Publisher confirms are a RabbitMQ extension to implement reliable publishing. When publisher confirms are enabled on
a channel, messages the client publishes are confirmed asynchronously by the broker, meaning they have been taken 
care of on the server side. We can imagine publisher confirms as acknowledges that the server send to the producer.
Publisher confirms are not enabled by default and must be enabled at the channel level. In the example two types
of approaches are explored. 

![07 diagram](./assets/07_new.png)

In the first one we send some messages and wait synchronously the broker confirmation. The confirmation is performed
serially, that is, after each message we wait the confirmation (so we don't batch published messages). In the second 
approach we send messages in batches and similarly we wait for confirmations in batches. After every n messages sent 
we want to receive the broker confirmation for all of them. Waiting for a batch of messages to be confirmed improves
throughput drastically over waiting for individual messages confirmations. A third approach could be used (not present
in the repo at the moment). The second approach has a better throughput, but still, we block the entire process waiting
for confirmations. Even worse, if one of the lasts messages in a batch are particularly slow, we are basically stopped
waiting a single message, with a available resources on both the client and the server side to process more concurrent
messages. So in the third approach we can asynchronously send messages and receive confirmations, with the process 
being completely independent. If we accumulate too many un-acked messages we could slow the producer process, but in
any case there isn't any drawback from slow confirmations from the other side of this concurrent system.

To start the example:
```shell
# start the producer with publisher confirms
go run ./07_pub-confirm
```
