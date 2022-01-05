package app

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

//Client of nats streaming server
type Client struct {
	url        string
	nc         *nats.Conn
	js         nats.JetStreamContext
	token      string
	subj       string
	streamName string
}

//New default client
func New(addr, tkn, subject, str string) *Client {
	return &Client{url: addr, token: tkn, subj: subject, streamName: str}
}

//Connect - init new nats client for subscribe
func (c *Client) connect(ctx context.Context) (err error) {

	if c.nc != nil && c.nc.IsConnected() {
		return
	}

	// Setup the connect options
	opts := []nats.Option{nats.Name("QueueService")}
	// Use UserCredentials

	if c.token != "" {
		opts = append(opts, nats.Token(c.token))
	}
	opts = append(opts, nats.DontRandomize())
	opts = append(opts, nats.ReconnectWait(5*time.Second))
	opts = append(opts, nats.MaxReconnects(-1))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Println("error", "nats.go DisconnectErrHandler", err)
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Println("info", "nats.go Reconnected", nc.ConnectedUrl())
	}))

	c.nc, err = nats.Connect(c.url, opts...)
	if err != nil {
		return
	}
	c.js, err = c.nc.JetStream()
	if err != nil {
		return
	}

	strInfo, err := c.js.StreamInfo(c.streamName)
	if err != nil || strInfo == nil {
		// Create the stream, which stores messages received on the subject.
		cfg := &nats.StreamConfig{
			Name:     c.streamName,
			Subjects: []string{c.subj},
			Storage:  nats.FileStorage,
			MaxAge:   15 * time.Minute,
		}

		if _, err = c.js.AddStream(cfg); err != nil {
			return
		}
	}

	c.getDelayed(ctx)
	//sub for enqueued messages
	_, err = c.js.QueueSubscribe(c.subj, "QueueService", c.queueMsg)
	c.worker(ctx)
	return
}

//Run streaming
func (c *Client) Run(ctx context.Context) {
	err := c.connect(ctx)
	if err != nil {
		log.Println("nats connect", err)
	}

	tk := time.NewTicker(15 * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			if c.nc != nil {
				c.nc.Drain()
			}
			log.Println("Queue Service shutdown")
			return
		case <-tk.C:
			err := c.connect(ctx)
			if err != nil {
				log.Println("nats connect", err)
			}
		}
	}
}
