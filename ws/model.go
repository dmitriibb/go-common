package ws

const (
	MessageTypeString MessageType = "string"
	MessageTypeLogin  MessageType = "login"
)

type MessageType string

type Message struct {
	Type    MessageType `json:"type"`
	Payload string      `json:"payload"`
}

type MessageLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
