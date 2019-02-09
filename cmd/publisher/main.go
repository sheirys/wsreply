package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sacOO7/gowebsocket"
	"github.com/sheirys/wsreply/broker"
	"github.com/sirupsen/logrus"
)

var kills = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}

type publisherState int

const (
	StateWaitingForSubscribers publisherState = iota
	StatePublishing
)

func main() {

	target := flag.String("c", "ws://localhost:8886/pub", "target url")
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, kills...)

	received := make(chan broker.Message, 5)
	send := make(chan broker.Message, 5)

	ticker := &time.Ticker{}
	state := StateWaitingForSubscribers

	ws := gowebsocket.New(*target)

	ws.OnConnectError = func(err error, ws gowebsocket.Socket) {
		logrus.WithError(err).WithField("target", *target).Fatal("connection failed")
	}

	ws.OnConnected = func(ws gowebsocket.Socket) {
		logrus.WithField("target", *target).Info("connected")
		send <- broker.MsgSyncSubscribers()
	}

	ws.OnTextMessage = func(payload string, ws gowebsocket.Socket) {
		var message broker.Message
		if err := json.Unmarshal([]byte(payload), &message); err != nil {
			logrus.WithError(err).WithField("payload", payload).Warn("cannot parse as broker.Message")
			return
		}
		logrus.WithFields(logrus.Fields{
			"op":   message.TranslateOp(),
			"data": string(message.Payload),
		}).Info("received message")

		received <- message
	}

	ws.OnDisconnected = func(err error, ws gowebsocket.Socket) {
		logrus.WithError(err).Error("disconnected")
		close(stop)
	}

	ws.Connect()

	for {
		select {
		case <-stop:
			ws.Close()
			return
		case msg := <-received:
			if msg.Op == broker.OpNoSubscribers {
				state = StateWaitingForSubscribers
				ticker.Stop()
				continue
			}

			if msg.Op == broker.OpHasSubscribers && state != StatePublishing {
				state = StatePublishing
				ticker = time.NewTicker(1 * time.Second)
				continue
			}
		case msg := <-send:
			// TODO: Make PR for ws.SendJSON(..)?
			logrus.WithFields(logrus.Fields{
				"op":   msg.TranslateOp(),
				"data": string(msg.Payload),
			}).Info("sending message")

			if bytes, err := json.Marshal(msg); err == nil {
				ws.SendBinary(bytes)
			}
		case <-ticker.C:
			send <- broker.MsgMessage([]byte("hello!"))
		}
	}
}
