# QS
Queue Service (QS) is an asynchronous message queuing service, distributed, horizontally scalable, persistent, time sorted message queue.

Standard queues provide at-least-once delivery, which means that each message is delivered at least once.

[Using Nats Jetstream.](https://nats.io)

Delay queues let you postpone the delivery of new messages to a queue for a number of seconds, for example, when your consumer application needs additional time to process messages. If you create a delay queue, any messages that you send to the queue remain invisible to consumers for the duration of the delay period. The default (minimum) delay for a queue is 0 seconds. The maximum is 15 minutes. 


Example

To delay message add to header time in seconds

pm.Header.Set("Delay-Time", "60")

Required parameters in the header

pm.Header.Set("Nats-Msg-Id", "unique id")
pm.Header.Set("Reply-Subject", "subject name to receive the queue message")


```

nc, err := nats.Connect(nats.DefaultURL, nats.Token("123wexsx2asekcijyc"))
	if err != nil {
		fmt.Println(err)
		return
	}

	js, err := nc.JetStream()
	if err != nil {
		fmt.Println(err)
		return
	}

	subj := "queueTest"
	strInfo, err := js.StreamInfo("streamName")
	if err != nil || strInfo == nil {
		// Create the stream, which stores messages received on the subject.
		cfg := &nats.StreamConfig{
			Name:     "streamName",
			Subjects: []string{subj},
			Storage:  nats.FileStorage,
			MaxAge:   15 * time.Minute,
		}

		if _, err = js.AddStream(cfg); err != nil {
			return
		}
	}

	_, err = js.QueueSubscribe(subj, "testGroup", func(m *nats.Msg) {
		fmt.Println(string(m.Data))
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	//without delay
	for i := 0; i < 1; i++ {
		id := strconv.FormatInt(time.Now().Unix(), 10)
		//default subject name
		pm := nats.NewMsg("EnQueue")
		//JetStream support idempotent message writes by ignoring
		//duplicate messages as indicated by the Nats-Msg-Id header.
		pm.Header.Set("Nats-Msg-Id", id+strconv.Itoa(i))
		pm.Header.Set("Reply-Subject", subj)
		pm.Data = []byte("message without delay, id: " + id + strconv.Itoa(i))
		_, err := js.PublishMsg(pm)
		if err != nil {
			fmt.Println("PublishMsg()", err)
			break
		}
	}

	//with delay
	for i := 25; i < 26; i++ {
		id := strconv.FormatInt(time.Now().Unix(), 10)
		//default subject name
		pm := nats.NewMsg("EnQueue")
		//JetStream support idempotent message writes by ignoring
		//duplicate messages as indicated by the Nats-Msg-Id header.
		pm.Header.Set("Nats-Msg-Id", id+strconv.Itoa(i))
		pm.Header.Set("Reply-Subject", subj)
		pm.Header.Set("Delay-Time", "60")
		pm.Data = []byte("message with delay, id: " + id + strconv.Itoa(i))
		_, err := js.PublishMsg(pm)
		if err != nil {
			fmt.Println("PublishMsg()", err)
			break
		}
	}

```