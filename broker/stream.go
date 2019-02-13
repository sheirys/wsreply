package broker

import "github.com/gorilla/websocket"

// Stream is used to save websocket connection in broker.
type Stream struct {
	isPublisher bool
	stream      *websocket.Conn
	broker      Broker
}
