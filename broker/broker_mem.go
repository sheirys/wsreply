package broker

import (
	"context"
	"log"
	"os"
)

type InMemBroker struct {
	ctx    context.Context
	queue  chan Message
	logger *log.Logger

	subscribers            map[*Stream]struct{}
	subscribe, unsubscribe chan *Stream
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
	b.logger.Printf("publishing message %#v", msg)
	b.queue <- msg
	return nil
}

func (b *InMemBroker) Start(ctx context.Context) error {
	b.ctx = ctx

	go func() {
		for {
			select {
			case <-b.ctx.Done():
				// TODO: call mass unsubscribe
				return
			case sub := <-b.subscribe:
				b.subscribers[sub] = struct{}{}
				b.broadcastNewSubscriber()
				b.logger.Printf("new stream %x (pub: %t) (streams: %d)", &sub, sub.isPublisher, len(b.subscribers))
			case unsub := <-b.unsubscribe:
				delete(b.subscribers, unsub)
				if len(b.subscribers) == 0 {
					b.broadcastNoSubscribers()
				}
				b.logger.Printf("unsubscribed stream %x (streams: %d)", &unsub, len(b.subscribers))
			case msg := <-b.queue:
				b.broadcastToSubscribers(msg)
			}
		}
	}()

	return nil
}

func NewInMemBroker() *InMemBroker {
	// FIXME: here is a lot of hardcoded sizes. Pass by argument or const?
	return &InMemBroker{
		queue:       make(chan Message, 5),
		subscribe:   make(chan *Stream, 5),
		unsubscribe: make(chan *Stream, 5),
		subscribers: make(map[*Stream]struct{}),
		logger:      log.New(os.Stdout, "broker-", 1),
	}
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
	for p := range b.subscribers {
		if !p.isPublisher {
			continue
		}
		select {
		case p.stream <- msg:
		default:
		}
	}
}
