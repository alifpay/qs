# QS
Queue Service (QS) is a fully managed message queuing service that enables you to decouple and scale microservices, distributed systems, and serverless applications. 

Standard queues provide at-least-once delivery, which means that each message is delivered at least once.

Depends on Nats Jetstream.

Delay queues let you postpone the delivery of new messages to a queue for a number of seconds, for example, when your consumer application needs additional time to process messages. If you create a delay queue, any messages that you send to the queue remain invisible to consumers for the duration of the delay period. The default (minimum) delay for a queue is 0 seconds. The maximum is 15 minutes. 


