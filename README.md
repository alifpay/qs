# QS
Queue Service (QS) is an asynchronous message queuing service, distributed, horizontally scalable, persistent, time sorted message queue.

Standard queues provide at-least-once delivery, which means that each message is delivered at least once.

[Using Nats Jetstream.](https://nats.io)

Delay queues let you postpone the delivery of new messages to a queue for a number of seconds, for example, when your consumer application needs additional time to process messages. If you create a delay queue, any messages that you send to the queue remain invisible to consumers for the duration of the delay period. The default (minimum) delay for a queue is 0 seconds. The maximum is 15 minutes. 


