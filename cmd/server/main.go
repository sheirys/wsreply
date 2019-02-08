package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sheirys/wsreply"
	"github.com/sheirys/wsreply/broker"
)

var kills = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}

func main() {

	listen := flag.String("l", "localhost:8886", "listen on")
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, kills...)

	app := &wsreply.Application{
		Broker: broker.NewInMemBroker(),
		Addr:   *listen,
		Log:    log.New(os.Stdout, "server-", 1),
	}
	app.Start()

	<-stop
	app.Stop()
}
