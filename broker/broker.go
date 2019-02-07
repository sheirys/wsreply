package broker

import (
	"context"

	"github.com/gorilla/websocket"
)

// Broker defines publish/subscribe system. In general this is simple fan
// out system but with some rules. When new sub
type Broker interface {

	// Start should initialize all variables required by broker. This
	// should be called before usage of other broker functions. When
	// context dies, all subscribers should be unsubscribed so
	// publishers know that there is nothing left.
	Start(context.Context) error

	// Broadcast should send message to all subscribers in broker.
	Broadcast(Message) error

	// NewPublisher should return new publishers stream, where messages
	// produced by broker should appear e.g. in this stream notifications
	// about new subcscriber connection or if there is no other
	// subscribers left should be pushed.
	AttachPublisherStream(*websocket.Conn) error

	// NewSubscriberStream should return new subscribers stream where
	// messages pushed by publishers should appear.
	AttachSubscriberStream(*websocket.Conn) error

	// Unsubscribe should close passed stream.
	Deattach(*websocket.Conn) error
}
