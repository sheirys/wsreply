package broker

import "github.com/gorilla/websocket"

type Stream struct {
	isPublisher bool
	stream      *websocket.Conn
	broker      Broker
}
