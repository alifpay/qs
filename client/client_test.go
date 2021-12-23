package client

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestClient(t *testing.T) {
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

	_, err = js.Subscribe(subj, func(m *nats.Msg) {
		fmt.Println(string(m.Data))
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	//without delay
	for i := 0; i < 25; i++ {
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
	for i := 25; i < 50; i++ {
		id := strconv.FormatInt(time.Now().Unix(), 10)
		//default subject name
		pm := nats.NewMsg("EnQueue")
		//JetStream support idempotent message writes by ignoring
		//duplicate messages as indicated by the Nats-Msg-Id header.
		pm.Header.Set("Nats-Msg-Id", id+strconv.Itoa(i))
		pm.Header.Set("Reply-Subject", subj)
		pm.Header.Set("Delay-Time", "7")
		pm.Data = []byte("message with delay, id: " + id + strconv.Itoa(i))
		_, err := js.PublishMsg(pm)
		if err != nil {
			fmt.Println("PublishMsg()", err)
			break
		}
	}

	time.Sleep(15 * time.Second)
}
