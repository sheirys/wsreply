package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	server := flag.String("s", "localhost", "server")
	port := flag.Int("p", 8882, "port")
	flag.Parse()

	bindAddr := fmt.Sprintf("%s:%d", *server, *port)

	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(bindAddr, nil)
}

// wsHandler will handle /ws url and tries to change http connection to ws.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		conn *websocket.Conn
		err  error
	)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	if conn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}
	for {
		mType, body, err := conn.ReadMessage()

		// close session if client disconnected or message type is not text
		if err != nil || mType != websocket.TextMessage {
			return
		}
		if err = conn.WriteMessage(mType, translate(body)); err != nil {
			return
		}
	}
}

// translate will change "?" symbols to "!" in provided message.
func translate(m []byte) []byte {
	return []byte(strings.Replace(string(m), "?", "!", -1))
}
