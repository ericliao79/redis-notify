package redis_notify

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

type Client struct {
	rcli *redis.Client

	rmsg chan *redis.Message

	eventHandle func(*Message)
}

func NewClient(rcli *redis.Client) *Client {
	return &Client{
		rcli: rcli,
		rmsg: make(chan *redis.Message),
	}
}

func (c *Client) Run(channel ...string) {
	pubsub := c.rcli.Subscribe(channel...)

	go func() {
		for {
			msg, _ := pubsub.ReceiveMessage()
			c.rmsg <- msg
		}
	}()

	for {
		select {
		case rmsg := <-c.rmsg:
			var m Message
			_ = json.Unmarshal([]byte(rmsg.Payload), &m)
			c.eventHandle(&m)
		}
	}
}

func (c *Client) SetEventHandle(f func(*Message)) {
	c.eventHandle = f
}
