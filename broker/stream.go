package broker

type Stream struct {
	isPublisher bool
	stream      chan Message
	broker      Broker
}

func (s *Stream) Read() Message {
	return <-s.stream
}

func (s *Stream) Close() {
	s.broker.Unsubscribe(s)
}
