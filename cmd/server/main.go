package main

import (
	"log"
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
		Addr:   ":8886",
		Logger: log.New(os.Stdout, "server-", 1),
	}
	app.Start()

	<-stop
	app.Stop()
}
