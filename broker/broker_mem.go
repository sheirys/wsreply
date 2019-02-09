package broker

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type InMemBroker struct {
	Log   *logrus.Logger
	Debug bool
	wg    *sync.WaitGroup
	die   chan struct{}

	subscribers map[*websocket.Conn]bool
	unsubscribe chan *websocket.Conn
	subscribe   chan *Stream
	message     chan Message

	sync.RWMutex
}

func (b *InMemBroker) AttachSubscriberStream(ws *websocket.Conn) error {
	b.Log.WithField("connection", &ws).Info("new subscriber")
	b.subscribe <- &Stream{
		stream:      ws,
		isPublisher: false,
	}
	return nil
}

func (b *InMemBroker) AttachPublisherStream(ws *websocket.Conn) error {
	b.Log.WithField("connection", &ws).Info("new publisher")
	b.subscribe <- &Stream{
		stream:      ws,
		isPublisher: true,
	}
	return nil
}

func (b *InMemBroker) Deattach(ws *websocket.Conn) error {
	b.Log.WithField("connection", &ws).Info("disconnecting")
	b.unsubscribe <- ws
	return nil
}

func (b *InMemBroker) Broadcast(msg Message) error {
	b.Log.WithFields(logrus.Fields{
		"op":   msg.TranslateOp(),
		"data": string(msg.Payload),
	}).Info("streaming message")
	b.message <- msg
	return nil
}

func (b *InMemBroker) Start() error {
	if b.Debug {
		b.Log.SetLevel(logrus.DebugLevel)
	}
	b.Log.Info("starting broker")

	go func() {
		b.wg.Add(1)
		for {
			select {
			case <-b.die:
				b.dropAll()
				b.wg.Done()
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

	return nil
}

func (b *InMemBroker) Stop() error {
	close(b.die)
	b.wg.Wait()
	return nil
}

func NewInMemBroker(debug bool) *InMemBroker {
	// FIXME: here is a lot of hardcoded sizes. Pass by argument or const?
	return &InMemBroker{
		message:     make(chan Message, 5),
		subscribe:   make(chan *Stream, 5),
		unsubscribe: make(chan *websocket.Conn, 5),
		subscribers: make(map[*websocket.Conn]bool),
		die:         make(chan struct{}),
		wg:          &sync.WaitGroup{},
		Log:         logrus.New(),
		Debug:       debug,
	}
}

// broadcastNewSubscriber will notify all publishers that new subscribers has
// connected.
func (b *InMemBroker) broadcastNewSubscriber() {
	b.broadcastToPublishers(MsgNewSubscriber())
}

// broadcastNoSubscribers will notify all publishers that there is no
// subscribers left in broker.
func (b *InMemBroker) broadcastNoSubscribers() {
	b.broadcastToPublishers(MsgNoSubscribers())
}

func (b *InMemBroker) broadcastHasSubscribers() {
	b.broadcastToPublishers(MsgHasSubscribers())
}

// broadcastToSubscribers will broadcast message to all subscribers in broker.
func (b *InMemBroker) broadcastToSubscribers(msg Message) {
	b.Log.WithField("op", msg.TranslateOp()).Debug("broadcasting to subscribers")
	b.RLock()
	defer b.RUnlock()
	for s, isPublisher := range b.subscribers {
		if !isPublisher {
			go s.WriteJSON(msg)
		}
	}
}

// broadcastToPublishers will broadcast message to all publishers in broker.
func (b *InMemBroker) broadcastToPublishers(msg Message) {
	b.Log.WithField("op", msg.TranslateOp()).Debug("broadcasting to publishers")
	b.RLock()
	defer b.RUnlock()
	for p, isPublisher := range b.subscribers {
		if isPublisher {
			go p.WriteJSON(msg)
		}
	}
}

func (b *InMemBroker) handleSubscribe(s *Stream) {
	b.Lock()
	b.subscribers[s.stream] = s.isPublisher
	b.Unlock()
	if !s.isPublisher {
		b.broadcastHasSubscribers()
	}
}

func (b *InMemBroker) handleUnsubscribe(ws *websocket.Conn) {
	b.Lock()
	delete(b.subscribers, ws)
	b.Unlock()
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
	b.RLock()
	defer b.RUnlock()
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
	b.Log.WithField("count", len(b.subscribers)).Info("dropping connections")

	b.RLock()
	defer b.RUnlock()
	for conn := range b.subscribers {
		conn.Close()
	}
}
