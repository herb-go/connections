package connections

type ConnectionsConsumer interface {
	OnMessage(*Message)
	OnError(*Error)
	OnClose(ConnectionOutput)
	OnOpen(ConnectionOutput)
}

type EmptyConsumer struct {
}

func (e EmptyConsumer) OnMessage(*Message) {

}
func (e EmptyConsumer) OnError(*Error) {

}
func (e EmptyConsumer) OnClose(ConnectionOutput) {

}
func (e EmptyConsumer) OnOpen(ConnectionOutput) {

}

func Consume(i ConnectionsInput, c ConnectionsConsumer) {
	for {
		select {
		case m := <-i.Messages():
			go func() {
				c.OnMessage(m)
			}()
		case e := <-i.Errors():
			go func() {
				c.OnError(e)
			}()
		case conn := <-i.OnCloseEvents():
			go func() {
				c.OnClose(conn)
			}()
		case conn := <-i.OnOpenEvents():
			go func() {
				c.OnOpen(conn)
			}()
		}
	}
}
