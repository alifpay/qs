package client

import (
	"strconv"

	"github.com/nats-io/nats.go"
)

// Send - sends the message to the Queue Service
func EnQueue(js nats.JetStreamContext, queueSubject, replySubject, uniqueId string, delay int, data []byte) error {
	var err error
	m := nats.NewMsg(queueSubject)
	m.Header.Set("Reply-Subject", replySubject)

	if delay > 0 {
		m.Header.Set("Delay-Time", strconv.Itoa(delay))
	}

	if len(uniqueId) > 0 {
		m.Header.Set("Nats-Msg-Id", uniqueId)
	}

	m.Data = data

	if _, err = js.PublishMsg(m); err != nil {
		return err
	}
	return nil
}
