package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/sacOO7/gowebsocket"
	"github.com/sheirys/wsreply/broker"
)

var kills = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}

func main() {

	target := flag.String("c", "ws://localhost:8886/sub", "target url")
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, kills...)

	ws := gowebsocket.New(*target)

	ws.OnConnectError = func(err error, ws gowebsocket.Socket) {
		logrus.WithError(err).WithField("target", *target).Fatal("connection failed")
	}

	ws.OnConnected = func(ws gowebsocket.Socket) {
		logrus.WithField("target", *target).Info("connected")
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
	}

	ws.OnDisconnected = func(err error, ws gowebsocket.Socket) {
		logrus.WithError(err).Error("disconnected")
	}

	ws.Connect()

	<-stop
	ws.Close()
}
