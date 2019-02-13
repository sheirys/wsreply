package broker

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// InMemBroker is Broker implementation that satisfies broker.Broker interface.
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

// AttachSubscriberStream attach subscribers stream to provided websocket
// connection. Here messages received from publishers will be streamed.
func (b *InMemBroker) AttachSubscriberStream(ws *websocket.Conn) error {
	b.Log.WithField("connection", &ws).Info("new subscriber")
	b.subscribe <- &Stream{
		stream:      ws,
		isPublisher: false,
	}
	return nil
}

// AttachPublisherStream attach publisher stream to provided websocket
// connection. Here messages about subscriber connections will be streamed.
func (b *InMemBroker) AttachPublisherStream(ws *websocket.Conn) error {
	b.Log.WithField("connection", &ws).Info("new publisher")
	b.subscribe <- &Stream{
		stream:      ws,
		isPublisher: true,
	}
	return nil
}

// Deattach this websocket connection from all streams.
func (b *InMemBroker) Deattach(ws *websocket.Conn) error {
	b.Log.WithField("connection", &ws).Info("disconnecting")
	b.unsubscribe <- ws
	return nil
}

// Broadcast should send message to all subscribers in broker.
func (b *InMemBroker) Broadcast(msg Message) error {
	b.Log.WithFields(logrus.Fields{
		"op":   msg.TranslateOp(),
		"data": msg.Payload,
	}).Info("streaming message")
	b.message <- msg
	return nil
}

// Start initialize all variables required by broker. This should be called
// before use of other broker functions.
func (b *InMemBroker) Start() error {
	if b.Debug {
		b.Log.SetLevel(logrus.DebugLevel)
	}
	b.Log.Info("starting broker")

	b.wg.Add(1)
	go func() {
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

// Stop broker, disconnect all subscribers and publishers.
func (b *InMemBroker) Stop() error {
	close(b.die)
	b.wg.Wait()
	return nil
}

// NewInMemBroker returns InMemBroker with predefined configurations.
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

// broadcastNoSubscribers will notify all publishers that there is no
// subscribers left in broker.
func (b *InMemBroker) broadcastNoSubscribers() {
	b.broadcastToPublishers(MsgNoSubscribers())
}

// broadcastHasSubscribers notifies all publishers that there is some listening
// subscribers on broker.
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
			s.WriteJSON(msg)
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
			p.WriteJSON(msg)
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
