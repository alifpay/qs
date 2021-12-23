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
			sb, err = c.js.QueueSubscribeSync(c.subj, "delayQueueService")
			if err != nil {
				log.Println("getDelayed sub", err)
				continue
			}
			msgNum, _, err := sb.Pending()
			if err != nil {
				log.Println("getDelayed sub", err)
				sb.Unsubscribe()
			}
			for i := 0; i < msgNum; i++ {
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
