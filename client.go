package redis_notify

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"os"
)

type Client struct {
	rcli *redis.Client

	rmsg chan *redis.Message

	eventHandle func(*Message)
}

type Message struct {
	Channel string `json:"channel"`
	Event string `json:"event"`
	Message interface{} `json:"message"`
}

func NewClient() *Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return &Client{
		rcli: client,
		rmsg: make(chan *redis.Message),
	}
}

func (c *Client) Run(channel... string) {
	pubsub := c.rcli.Subscribe(channel...)

	go func() {
		for {
			msg, _ := pubsub.ReceiveMessage()
			c.rmsg <-msg
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
