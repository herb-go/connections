package message

//Adapter connection message adapter
type Adapter map[string]func(*Message) error

//Register register message handler  by given type.
func (a *Adapter) Register(msgtype string, handler func(*Message) error) {
	(*a)[msgtype] = handler
}

//Exec exec message.
//Return if handler by message type exists and any error rasied.
func (a Adapter) Exec(msg *Message) (bool, error) {
	handler, ok := a[msg.Type]
	if ok == false {
		return false, nil
	}
	err := handler(msg)
	if err != nil {
		return true, err
	}
	return true, nil
}

// NewAdapter create new message adapter
func NewAdapter() *Adapter {
	return &Adapter{}
}
