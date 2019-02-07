package broker

import "github.com/gorilla/websocket"

type Stream struct {
	isPublisher bool
	stream      *websocket.Conn
	broker      Broker
}

/*
func (s *Stream) Read() Message {
	return <-s.stream
}

func (s *Stream) ReadWithNotify() <-chan Message {
	return s.stream
}

func (s *Stream) Close() {
	s.broker.Unsubscribe(s)
}
*/
