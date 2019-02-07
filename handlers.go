package wsreply

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sheirys/wsreply/broker"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (a *Application) PublisherWS(w http.ResponseWriter, r *http.Request) {
	var (
		ws     *websocket.Conn
		stream *broker.Stream
		err    error
	)

	if ws, err = upgrader.Upgrade(w, r, nil); err != nil {
		a.Logger.Println(err)
		return
	}
	defer ws.Close()

	if stream, err = a.Broker.NewPublisherStream(); err != nil {
		a.Logger.Println(err)
		return
	}
	defer stream.Close()

	for {
		_, body, err := ws.ReadMessage()
		if err != nil {
			a.Logger.Println(err)
			return
		}

		if err = ws.WriteMessage(websocket.TextMessage, body); err != nil {
			a.Logger.Println(err)
			return
		}

	}
}

func (a *Application) SubscriberWS(w http.ResponseWriter, r *http.Request) {
	var (
		ws     *websocket.Conn
		stream *broker.Stream
		err    error
	)

	if ws, err = upgrader.Upgrade(w, r, nil); err != nil {
		a.Logger.Println(err)
		return
	}
	defer ws.Close()

	if stream, err = a.Broker.NewSubscriberStream(); err != nil {
		a.Logger.Println(err)
		return
	}
	defer stream.Close()

	for {
		select {
		case message := <-stream.ReadWithNotify():
			if err = ws.WriteMessage(websocket.TextMessage, message.Payload); err != nil {
				a.Logger.Println(err)
				return
			}
		}

	}
}
