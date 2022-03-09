package app

import (
	"log"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

// to process enqueued message
func (c *Client) queueMsg(m *nats.Msg) {

	go func() {
		//parse header
		mt, err := m.Metadata()
		if err != nil {
			log.Println("m.Metadata()", err)
			return
		}
		replyTo, _, delayTime := c.parseHeader(m.Header, mt.Sequence.Stream)
		if delayTime == 0 || checkTime(mt.Timestamp) {
			pm := nats.NewMsg(replyTo)
			pm.Data = make([]byte, len(m.Data))
			copy(pm.Data, m.Data)
			_, err := c.js.PublishMsg(pm)
			if err != nil {
				log.Println("PublishMsg()", err)
				return
			}
			c.del(mt.Sequence.Stream)
			return
		}
		//cache delayed message
		delayTime = correctTime(mt.Timestamp).Unix() + delayTime
		sdt := strconv.FormatInt(delayTime, 10) + strconv.FormatUint(mt.Sequence.Stream, 10)
		addCache(sdt)
	}()
}

func (c *Client) parseHeader(h nats.Header, seq uint64) (replyTo, msgid string, delayTime int64) {
	//subject name for repsonse
	replyTo = h.Get("Reply-Subject")
	if replyTo == "" {
		log.Println("Reply-Subject is empty")
		c.del(seq)
		return
	}

	msgid = h.Get("Msg-Id")

	//how long to delay the message in seconds
	//min 1 second, max 900 seconds
	//delay time is empty then publish immediately
	delayTime, _ = strconv.ParseInt(h.Get("Delay-Time"), 10, 32)
	return
}

func (c *Client) del(seq uint64) {
	err := c.js.DeleteMsg(c.streamName, seq)
	if err != nil {
		log.Println("DeleteMsg()", err)
	}
}

//check, is message too old
func checkTime(t time.Time) bool {
	t = correctTime(t)
	return time.Since(t) > 1*time.Second
}

//change to local timezone
func correctTime(t time.Time) time.Time {
	loc := time.Now().Location()
	return t.In(loc)
}
