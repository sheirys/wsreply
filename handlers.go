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
		ws  *websocket.Conn
		err error
	)

	if ws, err = upgrader.Upgrade(w, r, nil); err != nil {
		a.Logger.Println(err)
		return
	}
	defer ws.Close()

	if err := a.Broker.AttachPublisherStream(ws); err != nil {
		a.Logger.Println(err)
		return
	}
	defer a.Broker.Deattach(ws)

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

		a.Broker.Broadcast(broker.Message{
			Op:      broker.OpMessage,
			Payload: Translate(body),
		})

	}
}

func (a *Application) SubscriberWS(w http.ResponseWriter, r *http.Request) {
	var (
		ws  *websocket.Conn
		err error
	)

	if ws, err = upgrader.Upgrade(w, r, nil); err != nil {
		a.Logger.Println(err)
		return
	}
	defer ws.Close()

	if err = a.Broker.AttachSubscriberStream(ws); err != nil {
		a.Logger.Println(err)
		return
	}
	defer a.Broker.Deattach(ws)

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
