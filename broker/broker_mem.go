package broker

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type InMemBroker struct {
	Log *log.Logger
	ctx context.Context
	wg  *sync.WaitGroup

	subscribers map[*websocket.Conn]bool
	unsubscribe chan *websocket.Conn
	subscribe   chan *Stream
	message     chan Message
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
	b.Log.Printf("broadcasting message %#v", msg)
	b.message <- msg
	return nil
}

func (b *InMemBroker) Start(ctx context.Context, wg *sync.WaitGroup) error {
	b.ctx = ctx
	b.wg = wg

	b.wg.Add(1)

	go func() {
		for {
			select {
			case <-b.ctx.Done():
				b.dropAll()
				return
			case sub := <-b.subscribe:
				b.handleSubscribe(sub)
			case unsub := <-b.unsubscribe:
				b.handleUnsubscribe(unsub)
			case msg := <-b.message:
				b.handleMessage(msg)
			}
		}
	}()

	b.wg.Done()
	return nil
}

func NewInMemBroker() *InMemBroker {
	// FIXME: here is a lot of hardcoded sizes. Pass by argument or const?
	return &InMemBroker{
		message:     make(chan Message, 5),
		subscribe:   make(chan *Stream, 5),
		unsubscribe: make(chan *websocket.Conn, 5),
		subscribers: make(map[*websocket.Conn]bool),
		Log:         log.New(os.Stdout, "broker-", 1),
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

func (b *InMemBroker) broadcastHasSubscribers() {
	b.broadcastToPublishers(Message{
		Op: OpHasSubscribers,
	})
}

// broadcastToSubscribers will broadcast message to all subscribers in broker.
func (b *InMemBroker) broadcastToSubscribers(msg Message) {
	for s, isPublisher := range b.subscribers {
		if !isPublisher {
			go s.WriteJSON(msg)
		}
	}
}

// broadcastToPublishers will broadcast message to all publishers in broker.
func (b *InMemBroker) broadcastToPublishers(msg Message) {
	for p, isPublisher := range b.subscribers {
		if isPublisher {
			go p.WriteJSON(msg)
		}
	}
}

func (b *InMemBroker) handleSubscribe(s *Stream) {
	b.subscribers[s.stream] = s.isPublisher
	if !s.isPublisher {
		b.broadcastHasSubscribers()
	}
	b.Log.Printf("new stream %x (pub: %t) (streams: %d)", &s, s.isPublisher, len(b.subscribers))
}

func (b *InMemBroker) handleUnsubscribe(ws *websocket.Conn) {
	delete(b.subscribers, ws)
	b.Log.Printf("unsubscribed stream %x (streams: %d)", &ws, len(b.subscribers))
	b.handleOpSyncSubscribers()
}

func (b *InMemBroker) handleMessage(m Message) {
	switch m.Op {
	case OpMessage:
		b.broadcastToSubscribers(m)
	case OpSyncSubscribers:
		b.handleOpSyncSubscribers()
	case OpNoSubscribers:
		b.broadcastNoSubscribers()
	}
}

func (b *InMemBroker) handleOpSyncSubscribers() {
	for _, isPublisher := range b.subscribers {
		if !isPublisher {
			b.broadcastHasSubscribers()
			return
		}
	}
	b.broadcastNoSubscribers()
	return
}

func (b *InMemBroker) dropAll() {
	b.Log.Println("dropping connections")
	for conn := range b.subscribers {
		conn.Close()
	}
	b.Log.Println("done")
}
