package broker

import (
	"github.com/gorilla/websocket"
)

// Broker defines publish/subscribe system. In general this is simple fan
// out system but with some additional rules.
type Broker interface {

	// Start should initialize all variables required by broker. This
	// should be called before use of other broker functions.
	Start() error

	// Stop should stop broker, disconnect all subscribers and publishers.
	Stop() error

	// Broadcast should send message to all subscribers in broker.
	Broadcast(Message) error

	// AttachPublisherStream should attach publisher stream to provided
	// websocket connection. Here messages about subscriber connections
	// should be streamed.
	AttachPublisherStream(*websocket.Conn) error

	// AttachSubscriberStream should attach subscribers stream to provided
	// websocket connection. Here messages received from publishers should
	// be streamed.
	AttachSubscriberStream(*websocket.Conn) error

	// Deattach this websocket connection from all streams.
	Deattach(*websocket.Conn) error
}
