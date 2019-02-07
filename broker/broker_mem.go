package broker

type InMemBroker struct {
	queue chan Message
	die   chan struct{}

	subscribers, publishers map[*Stream]struct{}
	subscribe, unsubscribe  chan *Stream
}

func (b *InMemBroker) NewSubscriberStream() (*Stream, error) {
	s := &Stream{
		stream:      make(chan Message, 5),
		isPublisher: false,
		broker:      b,
	}
	b.subscribe <- s
	return s, nil
}

func (b *InMemBroker) NewPublisherStream() (*Stream, error) {
	s := &Stream{
		stream:      make(chan Message, 5),
		isPublisher: true,
		broker:      b,
	}
	b.subscribe <- s
	return s, nil
}

func (b *InMemBroker) Unsubscribe(s *Stream) error {
	b.unsubscribe <- s
	return nil
}

func (b *InMemBroker) Publish(msg Message) error {
	b.queue <- msg
	return nil
}

func (b *InMemBroker) Stop() error {
	close(b.die)
	return nil
}

func (b *InMemBroker) Start() error {
	go func() {
		for {
			select {
			case <-b.die:
				return
			case sub := <-b.subscribe:
				b.subscribers[sub] = struct{}{}
				b.broadcastNewSubscriber()
			case unsub := <-b.unsubscribe:
				delete(b.subscribers, unsub)
				if len(b.subscribers) == 0 {
					b.broadcastNoSubscribers()
				}
			case msg := <-b.queue:
				b.broadcastToSubscribers(msg)
			}
		}
	}()

	return nil
}

// broadcastNewSubscriber will notify all publishers that new subscribers has
// connected.
func (b *InMemBroker) broadcastNewSubscriber() {
	b.broadcastToPublishers(Message{
		Op: OpNewSubscriber,
	})
}

// broadcastNoSubscribers will notify all publishers that there is no
// subscribers left in broker.
func (b *InMemBroker) broadcastNoSubscribers() {
	b.broadcastToPublishers(Message{
		Op: OpNoSubscribers,
	})
}

// broadcastToSubscribers will broadcast message to all subscribers in broker.
func (b *InMemBroker) broadcastToSubscribers(msg Message) {
	for s := range b.subscribers {
		if s.isPublisher {
			continue
		}
		select {
		case s.stream <- msg:
		default:
		}
	}
}

// broadcastToPublishers will broadcast message to all publishers in broker.
func (b *InMemBroker) broadcastToPublishers(msg Message) {
	for p := range b.publishers {
		if !p.isPublisher {
			continue
		}
		select {
		case p.stream <- msg:
		default:
		}
	}
}
