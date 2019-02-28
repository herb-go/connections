package message

type Message struct {
	Type string
	Room string
	Data interface{}
}

func New() *Message {
	return &Message{}
}
