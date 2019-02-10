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

func (s *Server) router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/sub", s.WSSubscriber)
	mux.HandleFunc("/pub", s.WSPublisher)

	return mux
}

// WSPublisher is endpoint fo publishers.
// Endpoint: /pub
func (s *Server) WSPublisher(w http.ResponseWriter, r *http.Request) {
	var (
		ws  *websocket.Conn
		err error
	)

	if ws, err = upgrader.Upgrade(w, r, nil); err != nil {
		s.Log.WithError(err).Warn("cannot upgrade connection")
		return
	}
	defer ws.Close()

	if err := s.Broker.AttachPublisherStream(ws); err != nil {
		s.Log.WithError(err).Warn("cannot attach publisher stream")
		return
	}
	defer s.Broker.Deattach(ws)

	for {
		message := broker.Message{}
		if err := ws.ReadJSON(&message); err != nil {
			s.Log.WithError(err).Warn("cannot unmarshal data")
			return
		}
		message.Payload = Translate(message.Payload)
		s.Broker.Broadcast(message)
	}
}

// WSSubscriber is endpoint for subscribers.
// Endpoint: /sub
func (s *Server) WSSubscriber(w http.ResponseWriter, r *http.Request) {
	var (
		ws  *websocket.Conn
		err error
	)

	if ws, err = upgrader.Upgrade(w, r, nil); err != nil {
		s.Log.WithError(err).Warn("cannot upgrade connection")
		return
	}
	defer ws.Close()

	if err = s.Broker.AttachSubscriberStream(ws); err != nil {
		s.Log.WithError(err).Warn("cannot attach subscriber stream")
		return
	}
	defer s.Broker.Deattach(ws)

	for {
		_, body, err := ws.ReadMessage()
		if err != nil {
			s.Log.WithError(err).Warn("connection error")
			return
		}

		if err = ws.WriteMessage(websocket.TextMessage, body); err != nil {
			s.Log.WithError(err).Warn("connection error")
			return
		}
	}
}
