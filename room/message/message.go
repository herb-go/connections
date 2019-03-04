package message

// Message connection message struct
type Message struct {
	//Type message type
	Type string
	// Room message room id
	Room string
	//Data message data
	Data interface{}
}

//New create new message.
func New() *Message {
	return &Message{}
}
