package ws

const (
	MessageTypeString MessageType = "string"
)

type MessageType string

type Message struct {
	Type    MessageType
	Payload string
}
