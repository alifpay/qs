package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func (c *Client) getDelayed(ctx context.Context) {
	var (
		sb  *nats.Subscription
		err error
	)
	for {
		select {
		case <-ctx.Done():
			log.Println("getDelayed is done")
			if sb != nil {
				sb.Unsubscribe()
			}
			return
		default:
			sb, err = c.js.QueueSubscribeSync(c.subj, "delayQueueService", nats.DeliverAll())
			if err != nil {
				log.Println("getDelayed sub", err)
				continue
			}
			si, err := c.js.StreamInfo(streamName)
			if err != nil {
				log.Println("getDelayed sub", err)
				sb.Unsubscribe()
			}

			var i uint64
			for i = 0; i < si.State.Msgs; i++ {
				m, err := sb.NextMsg(1 * time.Second)
				if err != nil {
					log.Println("getDelayed sub", err)
					sb.Unsubscribe()
					return
				}
				c.queueMsg(m)
				fmt.Println(string(m.Data))
			}

			sb.Unsubscribe()
			log.Println("getDelayed is done")
			return
		}
	}
}
