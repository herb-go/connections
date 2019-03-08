package connections

//Consumer Consumer interface which can consume connections message.
type Consumer interface {
	//OnMessage called when connection message received.
	OnMessage(*Message)
	//OnError called when onconnection error raised.
	OnError(*Error)
	//OnClose called when connection closed.
	OnClose(OutputConnection)
	//OnOpen called when connection open.
	OnOpen(OutputConnection)
	// Stop stop consumer
	Stop()
}

//EmptyConsumer empty consumer  implemented all  required method.
//Your consumer should extends to keep compatibility in future.
type EmptyConsumer struct {
}

//OnMessage called when connection message received.
func (e EmptyConsumer) OnMessage(*Message) {

}

//OnError called when onconnection error raised.
func (e EmptyConsumer) OnError(*Error) {

}

//OnClose called when connection closed.
func (e EmptyConsumer) OnClose(OutputConnection) {

}

//OnOpen called when connection open.
func (e EmptyConsumer) OnOpen(OutputConnection) {

}

// Stop stop consumer
func (e EmptyConsumer) Stop() {

}

// Consume consume input service with given consumer.
func Consume(i InputService, c Consumer) {
	for {
		select {
		case <-i.C():
			go func() {
				c.Stop()
			}()
			return
		case m := <-i.MessagesChan():
			go func() {
				c.OnMessage(m)
			}()
		case e := <-i.ErrorsChan():
			go func() {
				c.OnError(e)
			}()
		case conn := <-i.OnCloseEventsChan():
			go func() {
				c.OnClose(conn)
			}()
		case conn := <-i.OnOpenEventsChan():
			go func() {
				c.OnOpen(conn)
			}()
		}
	}
}
