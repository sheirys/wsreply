package broker

import (
	"context"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

type InMemBroker struct {
	ctx    context.Context
	queue  chan Message
	logger *log.Logger

	subscribers map[*websocket.Conn]bool
	unsubscribe chan *websocket.Conn
	subscribe   chan *Stream
}

func (b *InMemBroker) AttachSubscriberStream(ws *websocket.Conn) error {
	s := &Stream{
		stream:      ws,
		isPublisher: false,
	}
	b.subscribe <- s
	return nil
}

func (b *InMemBroker) AttachPublisherStream(ws *websocket.Conn) error {
	s := &Stream{
		stream:      ws,
		isPublisher: true,
	}
	b.subscribe <- s
	return nil
}

func (b *InMemBroker) Deattach(ws *websocket.Conn) error {
	b.unsubscribe <- ws
	return nil
}

func (b *InMemBroker) Broadcast(msg Message) error {
	b.logger.Printf("broadcasting message %#v", msg)
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
				b.subscribers[sub.stream] = sub.isPublisher
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
		unsubscribe: make(chan *websocket.Conn, 5),
		subscribers: make(map[*websocket.Conn]bool),
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
	for s, isPublisher := range b.subscribers {
		if isPublisher {
			continue
		}
		go s.WriteJSON(msg)
	}
}

// broadcastToPublishers will broadcast message to all publishers in broker.
func (b *InMemBroker) broadcastToPublishers(msg Message) {
	for p, isPublisher := range b.subscribers {
		if !isPublisher {
			continue
		}
		go p.WriteJSON(msg)
	}
}
