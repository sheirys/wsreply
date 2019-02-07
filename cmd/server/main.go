package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sheirys/wsreply"
	"github.com/sheirys/wsreply/broker"
)

var kills = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, kills...)

	app := &wsreply.Application{
		Broker: broker.NewInMemBroker(),
	}

	<-stop
	app.Stop()
}
