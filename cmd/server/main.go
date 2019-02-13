package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sheirys/wsreply"
	"github.com/sheirys/wsreply/broker"
	"github.com/sirupsen/logrus"
)

var kills = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}

func main() {

	listen := flag.String("l", "localhost:8886", "listen on")
	debug := flag.Bool("d", false, "debug mode")
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, kills...)

	srv := &wsreply.Server{
		Broker: broker.NewInMemBroker(*debug),
		Addr:   *listen,
		Log:    logrus.New(),
	}
	if err := srv.Init(); err != nil {
		logrus.WithError(err).Fatal("cannot init")
	}
	if err := srv.StartBroker(); err != nil {
		logrus.WithError(err).Fatal("cannot start broker")
	}
	if err := srv.StartHTTP(); err != nil {
		logrus.WithError(err).Fatal("cannot start http server")
	}

	<-stop

	if err := srv.Stop(); err != nil {
		logrus.WithError(err).Fatal("stopping server")
	}
}
