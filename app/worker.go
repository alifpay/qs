package app

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

//process  cached message every second
func (c *Client) worker(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			log.Println("worker is done")
			return
		default:
			time.Sleep(1 * time.Second)
			if qCache.unsorted {
				qCache.mu.Lock()
				qCache.items.Sort()
				qCache.unsorted = false
				qCache.mu.Unlock()
			}
			qCache.mu.Lock()
			now := time.Now().Unix()
			idx := -1
			for k, v := range qCache.items {
				tm, _ := strconv.ParseInt(v[:10], 10, 64)
				if now >= tm {
					go c.send(v, v[10:])
					idx = k
				} else {
					break
				}
			}
			if idx >= 0 && idx <= len(qCache.items) {
				qCache.items = qCache.items[idx+1:]
			}
			qCache.mu.Unlock()
		}
	}
}

func (c *Client) send(idx, seqStr string) {
	seq, _ := strconv.ParseUint(seqStr, 10, 64)
	m, err := c.js.GetMsg(streamName, seq)
	if err != nil {
		//todo err type if network retry later
		log.Println("c.js.GetMsg()", err)
		return
	}

	replyTo := m.Header.Get("Reply-Subject")
	pm := nats.NewMsg(replyTo)
	pm.Header.Set("Nats-Msg-Id", m.Header.Get("Nats-Msg-Id"))
	pm.Data = make([]byte, len(m.Data))
	copy(pm.Data, m.Data)
	_, err = c.js.PublishMsg(pm)
	if err != nil {
		log.Println("PublishMsg()", err)
		return
	}
	c.del(seq)
}
