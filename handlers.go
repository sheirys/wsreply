package wsreply

import (
	"encoding/json"
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
		a.Log.WithError(err).Warn("cannot upgrade connection")
		return
	}
	defer ws.Close()

	if err := a.Broker.AttachPublisherStream(ws); err != nil {
		a.Log.WithError(err).Warn("cannot attach publisher stream")
		return
	}
	defer a.Broker.Deattach(ws)

	for {
		_, body, err := ws.ReadMessage()
		if err != nil {
			a.Log.WithError(err).Warn("connection error")
			return
		}

		message := broker.Message{}
		if err := json.Unmarshal(body, &message); err != nil {
			a.Log.WithError(err).Warn("cannot unmarshal data")
			return
		}
		message.Payload = Translate(message.Payload)
		a.Broker.Broadcast(message)
	}
}

func (a *Application) SubscriberWS(w http.ResponseWriter, r *http.Request) {
	var (
		ws  *websocket.Conn
		err error
	)

	if ws, err = upgrader.Upgrade(w, r, nil); err != nil {
		a.Log.WithError(err).Warn("cannot upgrade connection")
		return
	}
	defer ws.Close()

	if err = a.Broker.AttachSubscriberStream(ws); err != nil {
		a.Log.WithError(err).Warn("cannot attach subscriber stream")
		return
	}
	defer a.Broker.Deattach(ws)

	for {
		_, body, err := ws.ReadMessage()
		if err != nil {
			a.Log.WithError(err).Warn("connection error")
			return
		}

		if err = ws.WriteMessage(websocket.TextMessage, body); err != nil {
			a.Log.WithError(err).Warn("connection error")
			return
		}
	}
}
