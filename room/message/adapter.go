package message

type Adapter map[string]func(*Message) error

func (a Adapter) Add(msgtype string, handler func(*Message) error) {
	a[msgtype] = handler
}
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

func NewAdapter() Adapter {
	return Adapter{}
}
