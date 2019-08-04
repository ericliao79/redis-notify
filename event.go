package redis_notify

type MessageType string

const (
	Msg       MessageType = "msg"
	Broadcast MessageType = "broadcast"
)

type Message struct {
	Type    MessageType `json:"type"`
	To      string      `json:"to"`
	Message interface{} `json:"message"`
}
